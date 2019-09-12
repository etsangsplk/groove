package rank

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/bbalet/stopwords"
	"github.com/biogo/ncbi/entrez"
	"github.com/hscells/cqr"
	"github.com/hscells/ghost"
	"github.com/hscells/groove/analysis"
	"github.com/hscells/groove/combinator"
	"github.com/hscells/groove/learning"
	"github.com/hscells/groove/pipeline"
	"github.com/hscells/groove/stats"
	"github.com/hscells/merging"
	"github.com/hscells/transmute"
	"github.com/hscells/transmute/fields"
	"github.com/hscells/trecresults"
	"github.com/reiver/go-porterstemmer"
	"gopkg.in/jdkato/prose.v2"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var cm = merging.CoordinationLevelMatching{
	Occurances: make(map[string]float64),
}

var writeMu sync.Mutex

// clf is the actual implementation of coordination level fusion. The exported function is simply a wrapper.
func clf(query pipeline.Query, posting *Posting, e stats.EntrezStatisticsSource, options CLFOptions) (trecresults.ResultList, error) {
	norm := merging.MinMaxNorm

	switch q := query.Query.(type) {
	case cqr.BooleanQuery:
		if strings.Contains(q.Operator, "adj") {
			q.Operator = "and"
		}
		lists := make([]merging.Items, len(q.Children))
		r := make([]trecresults.ResultList, len(q.Children))

		for i, c := range q.Children {
			var err error
			r[i], err = clf(pipeline.NewQuery(query.Name, query.Topic, c), posting, e, options)
			if err != nil {
				return nil, err
			}
		}

		fmt.Printf("merging %d lists - ", len(lists))

		for i, result := range r {
			fmt.Print(".")
			lists[i] = merging.FromTRECResults(result)
		}

		var items merging.Items
		if q.Operator == cqr.AND {
			items = merging.CombSUM{}.Merge(lists)
		} else {
			items = merging.CombMNZ{}.Merge(lists)
		}
		//norm.Init(items)
		//items = merging.Normalise(norm, items)
		n := float64(time.Now().Unix())

		sort.Slice(items, func(i, j int) bool {
			if items[i].Score != items[j].Score {
				return items[i].Score > items[j].Score
			}
			// Break ties by publication date.
			a := 1 - ((n - float64(posting.DocDates[hash(items[i].Id)]+1)) / n)
			b := 1 - ((n - float64(posting.DocDates[hash(items[j].Id)]+1)) / n)

			return a+items[i].Score > b+items[j].Score
		})

		list := items.TRECResults(query.Topic)
		for i := range list {
			list[i].Score++
			list[i].Rank = int64(i + 1)
		}

		fmt.Println("lists merged!")
		return list, nil
	case cqr.Keyword:
		scorers := []Scorer{
			&BM25Scorer{s: e, p: posting, B: 0.1, K1: 1.2},
			&LnL2Scorer{s: e, p: posting},
			&TitleAbstractScorer{s: e, p: posting},
			&PosScorer{s: e, p: posting},
			&SumIDFScorer{s: e, p: posting},
			&PubDateScorer{s: e, p: posting},
			&DocLenScorer{s: e, p: posting},
		}

		lists := make([]trecresults.ResultList, len(scorers))

		fmt.Printf("%s %v", q.QueryString, q.Fields)
		defer func() {
			fmt.Println("[√]")
		}()

		for i, scorer := range scorers {
			for pmid := range posting.DocLens {
				score := 1.0
				var err error
				switch q.Fields[0] {
				case fields.Title:
					score, err = scorer.Score(q.QueryString, pmid, "ti")
				case fields.Abstract:
					score, err = scorer.Score(q.QueryString, pmid, "ab")
				case fields.MeshHeadings, fields.MeSHTerms, fields.MeSHSubheading, fields.MeSHMajorTopic, fields.FloatingMeshHeadings:
					if exp, ok := q.Options[cqr.ExplodedString]; ok {
						if v, ok := exp.(bool); ok && q.QueryString[len(q.QueryString)-1] != '#' {
							if v {
								q.QueryString += "#"
							}
						}
					}
					score, err = scorer.Score(q.QueryString, pmid, "mh")
				case fields.TitleAbstract, fields.TextWord:
					score, err = scorer.Score(q.QueryString, pmid, "ti", "ab")
				case fields.AllFields:
					score, err = scorer.Score(q.QueryString, pmid, "ti", "ab", "mh")
				default:
					fmt.Printf("CANNOT SCORE QUERY STRING: %s\n", q.QueryString)
					score = 0.0
				}
				if err != nil {
					panic(err)
				}
				//n := float64(time.Now().Unix())
				//score *= 1 - ((n - float64(posting.DocDates[hash(pmid)]+1)) / n)
				lists[i] = append(lists[i], &trecresults.Result{
					Topic:     query.Topic,
					Iteration: "0",
					DocId:     pmid,
					Score:     score,
					RunName:   query.Topic,
				})
			}

			list := lists[i]
			R := merging.FromTRECResults(list)
			norm.Init(R)
			list = merging.Normalise(norm, R).TRECResults(query.Topic)

			sort.Slice(list, func(i, j int) bool {
				if list[i].Score != list[j].Score {
					return list[i].Score > list[j].Score
				}
				// Break ties by publication date.
				n := float64(time.Now().Unix())
				a := 1 - ((n - float64(posting.DocDates[hash(list[i].DocId)]+1)) / n)
				b := 1 - ((n - float64(posting.DocDates[hash(list[j].DocId)]+1)) / n)
				return a+list[i].Score > b+list[j].Score
			})

			for i := range list {
				list[i].Score++
				list[i].Rank = int64(i + 1)
			}

			lists[i] = list
			fmt.Print(".")
		}

		items := make([]merging.Items, len(lists))
		for i, list := range lists {
			items[i] = merging.FromTRECResults(list)
		}
		fmt.Print(len(items))

		if options.ScorePubMed {
			pmids := make([]string, len(posting.DocLens))
			i := 0
			for pmid := range posting.DocLens {
				pmids[i] = pmid
				i++
			}

			pmRes, err := scoreWithPubMed(pmids, q, query.Topic, e)
			if err != nil {
				return nil, err
			}
			fmt.Printf("$%d$", len(pmRes))

			if options.OnlyScorePubMed {
				norm.Init(merging.FromTRECResults(pmRes))
				return merging.Normalise(norm, merging.FromTRECResults(pmRes)).TRECResults(query.Topic), nil
			}
			items = append(items, merging.FromTRECResults(pmRes))
			fmt.Print(".")
			fmt.Print(len(items))
		}

		merger := merging.CombMNZ{}
		res := merger.Merge(items).TRECResults(query.Topic)
		norm.Init(merging.FromTRECResults(res))
		fusion := merging.Normalise(norm, merging.FromTRECResults(res)).TRECResults(query.Topic)

		if options.RetrievalModel {
			var rm trecresults.ResultList
			for _, row := range fusion {
				if row.Score >= 0.1 {
					rm = append(rm, row)
				}
			}
			return rm, nil
		}

		return fusion, nil
	}
	return nil, nil
}

func clm(query pipeline.Query, posting *Posting, e stats.EntrezStatisticsSource) (trecresults.ResultList, error) {
	//norm := merging.MinMaxNorm
	switch q := query.Query.(type) {
	case cqr.BooleanQuery:
		r := make([]trecresults.ResultList, len(q.Children))
		lists := make([]merging.Items, len(r))
		for i, child := range q.Children {
			var err error
			r[i], err = clm(pipeline.NewQuery(query.Name, query.Topic, child), posting, e)
			if err != nil {
				return nil, err
			}
		}

		fmt.Printf("merging %d lists", len(lists))

		for i, result := range r {
			lists[i] = merging.FromTRECResults(result)
		}

		for range lists {
			fmt.Print(".")
		}

		items := cm.Merge(lists)

		sort.Slice(items, func(i, j int) bool {
			return items[i].Score > items[j].Score
		})
		list := items.TRECResults(query.Topic)
		for i := range list {
			list[i].Rank = int64(i + 1)
		}

		fmt.Println("lists merged!")
		return list, nil
	case cqr.Keyword:
		scorer := SumIDFScorer{s: e, p: posting}
		fmt.Printf("%s %v", q.QueryString, q.Fields)
		var list trecresults.ResultList
		defer func() { fmt.Println(" [√]") }()

		for pmid := range posting.DocLens {
			var score float64
			var err error
			switch q.Fields[0] {
			case fields.Title:
				score, err = scorer.Score(q.QueryString, pmid, "ti")
			case fields.Abstract:
				score, err = scorer.Score(q.QueryString, pmid, "ab")
			case fields.MeshHeadings, fields.MeSHTerms, fields.MeSHSubheading, fields.MeSHMajorTopic, fields.FloatingMeshHeadings:
				score, err = scorer.Score(q.QueryString, pmid, "mh")
			case fields.TitleAbstract, fields.TextWord:
				score, err = scorer.Score(q.QueryString, pmid, "ti", "ab")
			case fields.AllFields:
				score, err = scorer.Score(q.QueryString, pmid, "ti", "ab", "mh")
			default:
				score = 0.0
			}
			if err != nil {
				return nil, err
			}
			if score == 0 {
				continue
			}
			list = append(list, &trecresults.Result{
				Topic:     query.Topic,
				Iteration: "0",
				DocId:     pmid,
				Score:     score,
				RunName:   query.Topic,
			})
		}

		//l := merging.FromTRECResults(lists[i])
		//norm.Init(l)
		//list := merging.Normalise(norm, l).TRECResults(query.Topic)

		sort.Slice(list, func(i, j int) bool {
			return list[i].Score > list[j].Score
		})

		for i := range list {
			list[i].Rank = int64(i + 1)
		}
		fmt.Printf(".")
		return list, nil

	}
	return nil, nil
}

func scoreWithPubMed(pmids []string, query cqr.CommonQueryRepresentation, topic string, e stats.EntrezStatisticsSource) (trecresults.ResultList, error) {
	//rCacher := combinator.NewFileQueryCache("./queries_cache_r")
	//nrCacher := combinator.NewFileQueryCache("./queries_cache_nr")

	rCacher, err := ghost.Open("./queries_cache_r", ghost.NewGobSchema(combinator.Documents{}), ghost.WithIndexCache(1e4))
	if err != nil {
		return nil, err
	}
	nrCacher, err := ghost.Open("./queries_cache_nr", ghost.NewGobSchema(combinator.Documents{}), ghost.WithIndexCache(1e4))
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})

	if topic == "CD009263" || topic == "CD010409" {
		pmids = pmids[:4000]
	}

	pmidKeywords := make([]cqr.CommonQueryRepresentation, len(pmids))
	for i, pmid := range pmids {
		pmidKeywords[i] = cqr.NewKeyword(pmid, "pmid")
	}

	tq := cqr.NewBooleanQuery("and", []cqr.CommonQueryRepresentation{
		query,
		cqr.NewBooleanQuery("or", pmidKeywords),
	})

	key := strconv.Itoa(int(hash(tq.String())))

	var r combinator.Documents
	var nr combinator.Documents
	if err := rCacher.Get(key, &r); err == nil && r != nil {
		err = nrCacher.Get(key, &nr)
		if err != nil {
			fmt.Println(err)
			goto research
		}
		results := make(trecresults.ResultList, len(r)+len(nr))
		for i := 0; i < len(r); i++ {
			results[i] = &trecresults.Result{
				Topic:   topic,
				DocId:   strconv.Itoa(int(r[i])),
				Rank:    int64(i) + 1,
				Score:   1 - (float64(i+1) / float64(len(r))),
				RunName: "pubmed",
			}
		}
		for i, j := len(r), 0; i < len(r)+len(nr); i++ {
			results[i] = &trecresults.Result{
				Topic:   topic,
				DocId:   strconv.Itoa(int(nr[j])),
				Rank:    int64(i) + 1,
				Score:   0,
				RunName: "pubmed",
			}
			j++
		}
		return results, nil
	}

research:
	q, err := transmute.CompileCqr2PubMed(tq)
	if err != nil {
		return nil, err
	}
	ranking, err := e.Search(q, func(p *entrez.Parameters) {
		p.Sort = "relevance"
	})
	if err != nil {
		fmt.Println(err)
		time.Sleep(5 * time.Second)
		goto research
		//return nil, err
	}
	rD := make(combinator.Documents, len(ranking))
	results := make(trecresults.ResultList, len(pmids))
	for i, pmid := range ranking {
		rD[i] = combinator.Document(pmid)
		s := strconv.Itoa(pmid)
		seen[s] = struct{}{}
		results[i] = &trecresults.Result{
			Topic:   topic,
			DocId:   s,
			Rank:    int64(i) + 1,
			Score:   1 - (float64(i+1) / float64(len(rD))),
			RunName: "pubmed",
		}
	}
	nrD := make(combinator.Documents, len(pmids)-len(ranking))
	for i, j := len(ranking), 0; i < len(pmids); i++ {
		for _, pmid := range pmids {
			if _, ok := seen[pmid]; !ok {
				p, err := strconv.Atoi(pmid)
				if err != nil {
					return nil, err
				}
				nrD[j] = combinator.Document(p)
				j++
				results[i] = &trecresults.Result{
					Topic:   topic,
					DocId:   pmid,
					Rank:    int64(i) + 1,
					Score:   0,
					RunName: "pubmed",
				}
				seen[pmid] = struct{}{}
				break
			}
		}
	}
	fmt.Println()
	fmt.Println(len(ranking), len(rD), len(nrD))
	fmt.Println()
	err = rCacher.Put(key, rD)
	if err != nil {
		return nil, err
	}

	err = nrCacher.Put(key, nrD)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func writeResults(list trecresults.ResultList, dir string) error {
	f, err := os.OpenFile(dir, os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, res := range list {
		_, err := f.WriteString(fmt.Sprintf("%s 0 %s %d %f %s\n", res.Topic, res.DocId, res.Rank, res.Score, res.RunName))
		if err != nil {
			return err
		}
	}
	return nil
}

func clfVariations(query cqr.CommonQueryRepresentation, topic string, idealPosting *Posting, e stats.EntrezStatisticsSource, options CLFOptions) error {
	candidates, err := learning.Variations(learning.CandidateQuery{
		TransformationID: -1,
		Topic:            topic,
		Query:            query,
		Chain:            nil,
	}, e, analysis.NewMemoryMeasurementExecutor(), nil,
		learning.NewLogicalOperatorTransformer(),
		learning.NewFieldRestrictionsTransformer(),
		learning.NewMeshParentTransformer(),
		learning.NewClauseRemovalTransformer())
	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)

	cd, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	indexPath := path.Join(cd, "groove_rank_variations")

	p := "tar18t2_variations"
	err = os.MkdirAll(path.Join(p, topic), 0777)
	if err != nil {
		return err
	}

	N, err := e.RetrievalSize(query)
	if err != nil {
		return err
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	for i, candidate := range candidates {
		fmt.Printf("[%s] variation %d/%d\n", topic, i+1, len(candidates))

		// String-ify the query.
		s, err := transmute.CompileCqr2PubMed(candidate.Query)
		if err != nil {
			return err
		}
	r:
		// Skip this candidate if it retrieves more than the original query.
		n, err := e.RetrievalSize(candidate.Query)
		if err != nil {
			fmt.Println(err)
			goto r
		}
		if n < N/2 || n > N*2 || n == 0 {
			fmt.Printf("skipping variation %d, retrieved no documents\n", i+1)
			fmt.Println(s)
			continue
		}
	s:
		// Obtain list of pmids.
		pmids, err := e.Search(s)
		if err != nil {
			fmt.Println(err)
			goto s
		}
		// Create posting list for query.
	f:
		posting, err := newPostingFromPMIDS(pmids, topic+"_"+strconv.Itoa(int(hash(s))), indexPath, e)
		if err != nil {
			fmt.Println(err)
			goto f
		}
		// Use fusion technique to rank retrieved results and write results to file.
		res, err := clf(pipeline.NewQuery(topic, topic, candidate.Query), posting, e, options)
		if err != nil {
			return err
		}
		err = writeResults(res, path.Join(p, topic, strconv.Itoa(int(hash(s)))+".res.retrieved"))
		if err != nil {
			return err
		}
		// Use fusion technique to rank only the relevant results and write to file.
		idealRes, err := clf(pipeline.NewQuery(topic, topic, candidate.Query), idealPosting, e, options)
		if err != nil {
			return err
		}
		err = writeResults(idealRes, path.Join(p, topic, strconv.Itoa(int(hash(s)))+".res.ideal"))
		if err != nil {
			return err
		}
		// Write the query to file for posterity.
		f, err := os.OpenFile(path.Join(p, topic, strconv.Itoa(int(hash(s)))+".qry"), os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		_, err = f.WriteString(s)
		if err != nil {
			return err
		}
		err = f.Close()
	}
	wg.Wait()
	return nil
}

func expandPhrases(phrases ...string) []string {
	var expansions []string
	for _, phrase := range phrases {
		phrase = strings.ReplaceAll(strings.ToLower(phrase), `"`, "")
		if strings.ContainsAny(phrase, "*") {
			terms := strings.Split(phrase, " ")
			for _, term := range terms {
				if term[len(term)-1] == '*' {
					var newExp []string
					for _, suf := range suffixes {
						newExp = append(newExp, strings.Replace(phrase, term, fmt.Sprintf("%s%s", strings.Replace(term, `*`, "", -1), suf), -1))
					}
					expansions = append(expansions, expandPhrases(newExp...)...)
				}
			}
		} else {
			expansions = append(expansions, phrase)
		}
	}
	return expansions
}

type CLFOptions struct {
	CLF            bool `json:"clf"`
	RankCLM        bool `json:"rank_clm"`
	RankCLF        bool `json:"rank_clf"`
	QueryExpansion bool `json:"query_expansion"`

	PubMedBaseline bool `json:"pubmed_baseline"`

	ScorePubMed     bool `json:"score_pubmed"`
	OnlyScorePubMed bool `json:"only_score_pubmed"`

	RetrievalModel bool    `json:"retrieval_model"`
	Cutoff         float64 `json:"cutoff"`

	CLFVariations bool   `json:"clf_variations"`
	PMIDS         string `json:"pmids"`
	Titles        string `json:"titles"`
}

// CLF performs coordination-level fusion given a query. It ranks the documents retrieved for a query according to ...TODO?
// This wrapper function performs some pre-processing steps before actually ranking the documents for the query.
func CLF(query pipeline.Query, cacher combinator.QueryCacher, e stats.EntrezStatisticsSource, options CLFOptions) (trecresults.ResultList, error) {
	cd, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	indexPath := path.Join(cd, "groove_rank")
	idealIndexPath := path.Join(cd, "groove_rank_ideal")

	var pmids []int
	b, err := ioutil.ReadFile(path.Join(options.PMIDS, query.Topic))
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(bytes.NewBuffer(b))
	for s.Scan() {
		pmid, err := strconv.Atoi(s.Text())
		if err != nil {
			return nil, err
		}
		pmids = append(pmids, pmid)
	}

	posting, err := newPostingFromPMIDS(pmids, query.Topic, indexPath, e)

	if options.QueryExpansion {
		c := query.Query.(cqr.BooleanQuery).Children
		c = append(c, cqr.NewBooleanQuery(cqr.AND, []cqr.CommonQueryRepresentation{
			cqr.NewKeyword("sensitivity", fields.TitleAbstract),
			cqr.NewKeyword("specificity", fields.TitleAbstract),
			cqr.NewKeyword("diagnos*", fields.TitleAbstract),
			cqr.NewKeyword("diagnosis", fields.TitleAbstract),
			cqr.NewKeyword("predictive", fields.TitleAbstract),
			cqr.NewKeyword("accuracy", fields.TitleAbstract),
		}))

		title, err := ioutil.ReadFile(path.Join(options.Titles, query.Topic))
		if err != nil {
			return nil, err
		}
		fmt.Println(string(title))

		titleParsed, err := prose.NewDocument(stopwords.CleanString(string(title), "en", false), prose.WithTagging(false), prose.WithExtraction(false), prose.WithSegmentation(false))
		if err != nil {
			return nil, err
		}
		var titleKeywords []cqr.CommonQueryRepresentation
		for _, tok := range titleParsed.Tokens() {
			v := porterstemmer.StemString(tok.Text)
			titleKeywords = append(titleKeywords, cqr.NewKeyword(v+"*", fields.TitleAbstract), cqr.NewKeyword(tok.Text, fields.TitleAbstract))
		}
		fmt.Println(titleKeywords)
		c = append(c, cqr.NewBooleanQuery(cqr.AND, titleKeywords))

		q := query.Query.(cqr.BooleanQuery)
		q.Children = c
		query.Query = q
	}

	if options.PubMedBaseline {
		sPmids := make([]string, len(pmids))
		for i, pmid := range pmids {
			sPmids[i] = strconv.Itoa(pmid)
		}
		fmt.Println(len(sPmids))
		return scoreWithPubMed(sPmids, query.Query, query.Topic, e)
	}

	if options.CLFVariations {
		f, err := os.OpenFile("/Users/s4558151/Repositories/tar/2018-TAR/Task2/Testing/qrels/qrel_abs_task2", os.O_RDONLY, 0664)
		if err != nil {
			return nil, err
		}
		qrels, err := trecresults.QrelsFromReader(f)
		if err != nil {
			return nil, err
		}
		rels := qrels.Qrels[query.Topic]
		var pmidsIdeal []int
		for _, rel := range rels {
			if rel.Score > 0 {
				i, err := strconv.Atoi(rel.DocId)
				if err != nil {
					return nil, err
				}
				pmidsIdeal = append(pmids, i)
			}
		}
		idealPosting, err := newPostingFromPMIDS(pmidsIdeal, query.Topic, idealIndexPath, e)
		if err != nil {
			return nil, err
		}

		if _, err := os.Stat("/Users/s4558151/go/src/github.com/hscells/groove/scripts/tar18t2_variations/" + query.Topic); os.IsNotExist(err) {
			res, err := e.Execute(query, e.SearchOptions())
			if err != nil {
				return nil, err
			}
			err = writeResults(res, path.Join("/Users/s4558151/go/src/github.com/hscells/groove/scripts/tar18t2_orig/"+query.Topic))
			if err != nil {
				return nil, err
			}
			return nil, clfVariations(query.Query, query.Topic, idealPosting, e, options)
		} else {
			fmt.Printf("skipping topic %s, already exists\n", query.Topic)
		}
	} else if options.RankCLF {
		results, err := clf(query, posting, e, options)
		if err != nil {
			return nil, err
		}
		return results, nil
	} else if options.RankCLM {
		results, err := clm(query, posting, e)
		if err != nil {
			return nil, err
		}
		return results, nil
	}

	return nil, err
}