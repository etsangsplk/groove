#!/usr/bin/env bash

RANK_CLM=false
RANK_CLF=true
QUERY_EXPANSION=false
SCORE_PUBMED=true
ONLY_SCORE_PUBMED=false
RETRIEVAL_MODEL=true
CLF_VARIATIONS=false
PMIDS=$1
TITLES=$2
QRELS=$3
RUN=$4
QUERY_PATH=$5

CUTOFFS=(0.05 0.075 0.1 0.125 0.15 0.175 0.2 0.25 0.3)

for cut in ${CUTOFFS[@]}; do
    echo "running ${RUN} with cutoff ${cut}"
    ./clf.sh ${RANK_CLM} ${RANK_CLF} ${QUERY_EXPANSION} ${SCORE_PUBMED} ${ONLY_SCORE_PUBMED} ${RETRIEVAL_MODEL} ${CLF_VARIATIONS} ${PMIDS} ${TITLES} ${cut} ${RUN}${cut}.run ${QRELS} ${QUERY_PATH}
done