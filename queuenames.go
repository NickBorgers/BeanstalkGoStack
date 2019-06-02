package main

import (
        "os"
)

var requestQueue = "pandaDataRequest" + os.Getenv("ENV_NAME_MODIFIER")
var healthDataQueue = "pandaHealthData" + os.Getenv("ENV_NAME_MODIFIER")
var healthAnalysisQueue = "pandaHealthAnalysis" + os.Getenv("ENV_NAME_MODIFIER")
