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
		rest.Get("/List", func(w rest.ResponseWriter, req *rest.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			// FunctionName := req.PathParam("function")
			// switch {
			// case FunctionName == "List":
			ReturnResult := dashboardlib.ListProjectSummaryColumn()
			w.WriteJson(&ReturnResult)
			// }
		}),
		rest.Get("/#project/#type/#value", func(w rest.ResponseWriter, req *rest.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			ProjectName := req.PathParam("project")
			TypeName := req.PathParam("type")
			dataName := req.PathParam("value")
			switch {

			case ProjectName == "TYGH":
				switch {

				case TypeName == "WeekToDate":
					ReturnResult := dashboardlib.TYGHDiffWeekDate(ProjectName, dataName)
					w.WriteJson(&ReturnResult)
				}
			case ProjectName == "BABY":
				switch {

				case TypeName == "WeekToDate":
					ReturnResult := dashboardlib.BABYDiffWeekDate(ProjectName, dataName)
					w.WriteJson(&ReturnResult)
				}
			case ProjectName == "MMHDRUG":
				switch {

				case TypeName == "WeekToDate":
					ReturnResult := dashboardlib.DiffWeekDate(ProjectName, dataName)
					w.WriteJson(&ReturnResult)
				}

			}

		}),
		rest.Get("/#project/#type", func(w rest.ResponseWriter, req *rest.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			ProjectName := req.PathParam("project")
			TypeName := req.PathParam("type")

			switch {

			case TypeName == "LastUpdate":
				ReturnResult := dashboardlib.LastUpdateWeek()
				w.WriteJson(&ReturnResult)
			case TypeName == "ProjectSummary":
				ReturnResult := dashboardlib.ProjectSummary(ProjectName)
				w.WriteJson(&ReturnResult)
			case TypeName == "Status":
				ReturnResult := dashboardlib.PieChart(ProjectName, TypeName)
				w.WriteJson(&ReturnResult)
			case TypeName == "Priority":
				ReturnResult := dashboardlib.PieChart(ProjectName, TypeName)
				w.WriteJson(&ReturnResult)
			case TypeName == "Resolution":
				ReturnResult := dashboardlib.PieChart(ProjectName, TypeName)
				w.WriteJson(&ReturnResult)

			case ProjectName == "TYGH":
				switch {

				case TypeName == "DiffWeek":
					ReturnResult := dashboardlib.TYGHDiffWeek(ProjectName)
					w.WriteJson(&ReturnResult)

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

				case TypeName == "DiffVersionSoFar":
					ReturnResult := dashboardlib.DiffVersionSoFarRemain(ProjectName)
					w.WriteJson(&ReturnResult)

				case TypeName == "WeekRemain":
					ReturnResult := dashboardlib.WeekRemain(ProjectName)
					w.WriteJson(&ReturnResult)

				case TypeName == "IssueTimeSpent":
					ReturnResult := dashboardlib.IssueTimespent(ProjectName)
					w.WriteJson(&ReturnResult)

				case TypeName == "DueDateRemain":
					ReturnResult := dashboardlib.DueDateRemain(ProjectName)
					w.WriteJson(&ReturnResult)
				}

			case ProjectName == "BABY":
				switch {
				case TypeName == "IssueTimeSpent":
					ReturnResult := dashboardlib.BABYIssueTimespent()
					w.WriteJson(&ReturnResult)
				case TypeName == "DueDateRemain":
					ReturnResult := dashboardlib.BABYDueDateRemain()
					w.WriteJson(&ReturnResult)
				case TypeName == "DiffVersionSoFar":
					ReturnResult := dashboardlib.DiffVersionSoFarRemain(ProjectName)
					w.WriteJson(&ReturnResult)
				case TypeName == "DiffWeek":
					ReturnResult := dashboardlib.BABYDiffWeek(ProjectName)
					w.WriteJson(&ReturnResult)

				case TypeName == "DiffDate":
					ReturnResultA, ReturnResultI := dashboardlib.BABYDiffDate()
					var BABYReturnResult []dashboardlib.JsonResult
					BABYReturnResult = append(BABYReturnResult, ReturnResultA, ReturnResultI)
					w.WriteJson(&BABYReturnResult)
				case TypeName == "WeekRemain":
					ReturnResult := dashboardlib.WeekRemain(ProjectName)
					w.WriteJson(&ReturnResult)
				case TypeName == "ProjectSummary":
					ReturnResult := dashboardlib.ProjectSummary(ProjectName)
					w.WriteJson(&ReturnResult)
				}
			case ProjectName == "IOS":
				switch {
				case TypeName == "DiffVersionSoFar":
					ReturnResult := dashboardlib.BABYIOSDiffVersionSoFarRemain(ProjectName)
					w.WriteJson(&ReturnResult)
					// case TypeName == "DiffDate":
					// 	ReturnResult := dashboardlib.BABYIOSDiffDate(ProjectName)
					// 	w.WriteJson(&ReturnResult)
				}
			case ProjectName == "MMHDRUG":
				switch {
				case TypeName == "DiffWeek":
					ReturnResult := dashboardlib.DiffWeek(ProjectName)
					w.WriteJson(&ReturnResult)

				case TypeName == "IssueTimeSpent":
					ReturnResult := dashboardlib.IssueTimespent(ProjectName)
					w.WriteJson(&ReturnResult)
				case TypeName == "DueDateRemain":
					ReturnResult := dashboardlib.DueDateRemain(ProjectName)
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

				}
			}

		}),
	)

	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	http.Handle("/dashboard/", http.StripPrefix("/dashboard/", http.FileServer(http.Dir("./dashboard/"))))

	log.Fatal(http.ListenAndServe(":7108", nil))
}
