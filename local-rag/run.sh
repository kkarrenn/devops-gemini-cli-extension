#!/bin/sh
export RAG_DB_PATH=$(pwd)/devops-rag.db
completed=60
while [[ "$completed" != "18742" ]];do
   completed=$(./local-rag 2>&1 | tee -a embedding.log | grep "Exporting database Knowledge base docs"|cut -d: -f4 | cut -d, -f1)
   echo Docs done $completed
   sleep 65
done
