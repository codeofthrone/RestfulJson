package main

import (
	// _ "../github.com/GO-SQL-Driver/MySQL"
	// "../github.com/pmylund/sortutil"
	"./DashBoardLib"
	// "database/sql"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	// "net"
	"net/http"
)

func checkerr(str string, err error) {
	if err != nil {
		log.Println(str)
		log.Fatal(err.Error())
		panic(err.Error())
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/#project/#type", func(w rest.ResponseWriter, req *rest.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			ProjectName := req.PathParam("project")
			TypeName := req.PathParam("type")
			switch {
			case TypeName == "DiffVersionAPP":
				ReturnResult, _ := dashboardlib.TYGHDiffVersion(ProjectName)
				w.WriteJson(&ReturnResult)
			case TypeName == "DiffVersionWEB":
				_, ReturnResult := dashboardlib.TYGHDiffVersion(ProjectName)
				w.WriteJson(&ReturnResult)
			case TypeName == "DiffVersionSoFarAPP":
				ReturnResult, _ := dashboardlib.TYGHDiffVersionSoFarRemain(ProjectName)
				w.WriteJson(&ReturnResult)
			case TypeName == "DiffVersionSoFarWEB":
				_, ReturnResult := dashboardlib.TYGHDiffVersionSoFarRemain(ProjectName)
				w.WriteJson(&ReturnResult)

			case TypeName == "TYGHDiffDate":
				ReturnResultAPP, ReturnResultWEB := dashboardlib.TYGHDiffDate(ProjectName)
				var ReturnResult []dashboardlib.JsonResult
				ReturnResult = append(ReturnResult, ReturnResultAPP, ReturnResultWEB)
				w.WriteJson(&ReturnResult)

			case TypeName == "DiffVersion":
				ReturnResult := dashboardlib.DiffVersion(ProjectName)
				w.WriteJson(&ReturnResult)
			case TypeName == "DiffVersionSoFar":
				ReturnResult := dashboardlib.DiffVersionSoFarRemain(ProjectName)
				w.WriteJson(&ReturnResult)

			case TypeName == "DiffDate":
				ReturnResult := dashboardlib.DiffDate(ProjectName)
				w.WriteJson(&ReturnResult)
			case TypeName == "WeekRemain":
				ReturnResult := dashboardlib.WeekRemain(ProjectName)
				w.WriteJson(&ReturnResult)
			case TypeName == "IssueTimeSpent":
				ReturnResult := dashboardlib.IssueTimespent(ProjectName)
				w.WriteJson(&ReturnResult)
			case TypeName == "IssueType":
				fallthrough
			case TypeName == "Priority":
				fallthrough
			case TypeName == "Resolution":
				fallthrough
			case TypeName == "Status":
				ReturnResult := dashboardlib.PieChart(ProjectName, TypeName)
				w.WriteJson(&ReturnResult)

			case TypeName == "LastUpdate":
				ReturnResult := dashboardlib.LastUpdateWeek()
				w.WriteJson(&ReturnResult)

			}

		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":7108", api.MakeHandler()))
}
