package dashboardlib

import (
	_ "github.com/GO-SQL-Driver/MySQL"
	"github.com/pmylund/sortutil"
	"database/sql"
	// "html/template"
	"log"
	"strconv"
	"strings"
)

type JiraIssueVersions struct {
	Sn   int
	Name string
}

type JiraIssuesStatusDiffVersion struct {
	Name          string `json:"Name"`
	CreateCounter int    `json:"CreateCounter"`
	CloseCounter  int    `json:"CloseCounter"`
	RemainCounter int    `json:"RemainCounter"`
}

type JiraIssuesDateRemain struct {
	Name          string `json:"Name"`
	CreateCounter int    `json:"CreateCounter"`
	CloseCounter  int    `json:"CloseCounter"`
	RemainCounter int    `json:"RemainCounter"`
}

type JiraIssuesPriority struct {
	Fatal  string `json:"Name"`
	Low    string `json:"Name"`
	Medium string `json:"Name"`
	Name   string `json:"Name"`
}

type JsonResult struct {
	Name          []string `json:"Name"`
	CreateCounter []int    `json:"CreateCounter"`
	CloseCounter  []int    `json:"CloseCounter"`
	RemainCounter []int    `json:"RemainCounter"`
}

type JsonResultStr struct {
	Name          []string `json:"Name"`
	CreateCounter []string `json:"CreateCounter"`
	CloseCounter  []string `json:"CloseCounter"`
	RemainCounter []string `json:"RemainCounter"`
}

type JsonResultPie struct {
	Name          []string  `json:"Name"`
	RemainCounter []float64 `json:"RemainCounter"`
}

type JsonResultTable struct {
	Name     string `json:"Key"`
	Summary  string `json:"Summary"`
	DiffDate int    `json:"date"`
}

func checkerr(str string, err error) {
	if err != nil {
		log.Println(str)
		log.Fatal(err.Error())
		panic(err.Error())
	}
}

func QuerySingle(query string) string {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	InsertType, err := db.Query(query)
	checkerr(query, err)
	InsertType.Next()
	var tmp string
	InsertType.Scan(&tmp)
	InsertType.Close()
	db.Close()
	return tmp
}

func WeekRemain(ProjectName string) JsonResult {
	var ReturnJson JsonResult
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	WeekCreateSQL := "SELECT  WEEK(`CreatedTime`)+1 as week , count(*)  FROM Issues  WHERE `Key` LIKE '" + ProjectName + "-%'  group by week"
	WeekRemains, err := db.Query(WeekCreateSQL)
	// log.Println(WeekCreateSQL)
	checkerr(WeekCreateSQL, err)
	var WeekReaminResult []JiraIssuesDateRemain
	WeekRemainCounter := 0
	for WeekRemains.Next() {
		var (
			WeekRemain       JiraIssuesDateRemain
			tmpCreateCounter int
			tmpName          string
		)
		WeekRemains.Scan(&tmpName, &tmpCreateCounter)
		ClosqSQL := "Select closecounter from ( SELECT count(*) as closecounter, WEEK(UpdatedTime)+1 as week FROM Issues " +
			" WHERE `Key` LIKE '" + ProjectName + "-%' AND Status='Closed'  group by week ) as a WHERE a.week='" + tmpName + "'"
		tmpCloseCounter := QuerySingle(ClosqSQL)

		WeekRemain.Name = tmpName
		WeekRemain.CreateCounter = tmpCreateCounter
		WeekRemain.CloseCounter, _ = strconv.Atoi(tmpCloseCounter)
		WeekRemainCounter = WeekRemainCounter + tmpCreateCounter - WeekRemain.CloseCounter
		WeekRemain.RemainCounter = WeekRemainCounter
		WeekReaminResult = append(WeekReaminResult, WeekRemain)
	}
	WeekRemains.Close()
	sortutil.AscByField(WeekReaminResult, "Name")
	for _, value := range WeekReaminResult {
		// log.Println(key, value)
		ReturnJson.Name = append(ReturnJson.Name, value.Name)
		ReturnJson.CreateCounter = append(ReturnJson.CreateCounter, value.CreateCounter)
		ReturnJson.CloseCounter = append(ReturnJson.CloseCounter, value.CloseCounter)
		ReturnJson.RemainCounter = append(ReturnJson.RemainCounter, value.RemainCounter)
	}
	db.Close()
	// log.Println(ReturnJson)
	return ReturnJson
}

func VersionData(ProjectName string) []JiraIssuesStatusDiffVersion {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	// ProjectName := "TYGH"
	versSQL := "SELECT Id,Name from ( " +
		"SELECT f.Id,f.Name,f.Project FROM Fixversions as f LEFT join Versions as v on f.Id = v.Id UNION DISTINCT " +
		"SELECT v.Id,v.Name,v.Project FROM Fixversions as f RIGHT join Versions as v on f.Id = v.Id) as a where  Project='" + ProjectName + "'  group by a.Id " +
		"ORDER BY `a`.`Id` ASC"
	Vers, err := db.Query(versSQL)
	checkerr(versSQL, err)
	var DiffVersionResult []JiraIssuesStatusDiffVersion
	for Vers.Next() {
		var (
			IssueVersions JiraIssuesStatusDiffVersion
			tmpVerId      string
			tmpVerName    string
		)
		Vers.Scan(&tmpVerId, &tmpVerName)
		IssueVersions.Name = tmpVerName
		DiffVersionResult = append(DiffVersionResult, IssueVersions)
	}
	Vers.Close()

	// // TODO no version
	sortutil.AscByField(DiffVersionResult, "Name")
	IssuesSQL := "SELECT Id,FixVersions,Status,Version FROM Issues WHERE `Key` like '" + ProjectName + "-%' AND Resolution not like '% Fix'"
	Issuess, err := db.Query(IssuesSQL)
	checkerr(IssuesSQL, err)
	for Issuess.Next() {
		var (
			tmpId             string
			tmpFixVer         string
			tmpVer            string
			tmpStatus         string
			CreateVersionName []string
		)
		Issuess.Scan(&tmpId, &tmpFixVer, &tmpStatus, &tmpVer)
		// log.Println(tmpId, tmpFixVer, tmpStatus, tmpVer)
		IssueCreateCounter := 0
		if tmpVer != "0" {
			VerDataSQL := " SELECT Data FROM `Version` WHERE `Id` = '" + tmpId + "' AND Enable ='1' "
			VerDatas, err := db.Query(VerDataSQL)
			checkerr(VerDataSQL, err)
			for VerDatas.Next() {
				var tmpData string
				VerDatas.Scan(&tmpData)
				VersNameSQL := " SELECT Name FROM `Versions` WHERE `Sn` = '" + tmpData + "' AND Project ='" + ProjectName + "' "
				VerName := QuerySingle(VersNameSQL)
				// log.Println("Create", VerName)
				CreateVersionName = append(CreateVersionName, VerName)
				for key, Value := range DiffVersionResult {
					if Value.Name == VerName {
						DiffVersionResult[key].CreateCounter++
						IssueCreateCounter++
					}
				}
			}
			VerDatas.Close()
		}
		// log.Println("Create counter", IssueCreateCounter)
		if tmpStatus == "Closed" {
			FixVerDataSQL := " SELECT Data FROM `Fixversion` WHERE `Id` = '" + tmpId + "' AND Enable ='1' "
			VerDatas, err := db.Query(FixVerDataSQL)
			checkerr(FixVerDataSQL, err)
			for VerDatas.Next() {
				var tmpData string
				VerDatas.Scan(&tmpData)
				VersNameSQL := " SELECT Name FROM `Fixversions` WHERE `Sn` = '" + tmpData + "' AND Project ='" + ProjectName + "' "
				VerName := QuerySingle(VersNameSQL)
				// log.Println("close", VerName)
				for key, Value := range DiffVersionResult {
					if Value.Name == VerName && IssueCreateCounter > 0 {
						DiffVersionResult[key].CloseCounter++
						IssueCreateCounter--
					}
				}
			}
		}
		// log.Println("Close counter", IssueCreateCounter, VersionRemainCounter)
	}
	Issuess.Close()
	db.Close()
	return DiffVersionResult
}

func TYGHDiffVersion(ProjectName string) (JsonResult, JsonResult) {
	DiffVersionResult := VersionData(ProjectName)
	var ReturnJsonAPP JsonResult
	var ReturnJsonWEB JsonResult
	APPRemain := 0
	WEBRemain := 0
	appflag := 0
	webflag := 0
	for _, Value := range DiffVersionResult {
		if strings.Contains(Value.Name, "APP") {
			if appflag == 0 {
				ReturnJsonAPP.RemainCounter = append(ReturnJsonAPP.RemainCounter, 0)
			} else {
				ReturnJsonAPP.RemainCounter = append(ReturnJsonAPP.RemainCounter, APPRemain)
			}
			APPRemain = APPRemain + Value.CreateCounter - Value.CloseCounter
			ReturnJsonAPP.Name = append(ReturnJsonAPP.Name, Value.Name)
			ReturnJsonAPP.CreateCounter = append(ReturnJsonAPP.CreateCounter, Value.CreateCounter)
			ReturnJsonAPP.CloseCounter = append(ReturnJsonAPP.CloseCounter, Value.CloseCounter)
			appflag++
		} else {
			if webflag == 0 {
				ReturnJsonWEB.RemainCounter = append(ReturnJsonWEB.RemainCounter, 0)
			} else {
				ReturnJsonWEB.RemainCounter = append(ReturnJsonWEB.RemainCounter, WEBRemain)
			}
			ReturnJsonWEB.Name = append(ReturnJsonWEB.Name, Value.Name)
			WEBRemain = WEBRemain + Value.CreateCounter - Value.CloseCounter
			ReturnJsonWEB.CreateCounter = append(ReturnJsonWEB.CreateCounter, Value.CreateCounter)
			ReturnJsonWEB.CloseCounter = append(ReturnJsonWEB.CloseCounter, Value.CloseCounter)
			webflag++
		}
	}
	return ReturnJsonAPP, ReturnJsonWEB
}

func DiffVersion(ProjectName string) JsonResult {
	DiffVersionResult := VersionData(ProjectName)
	var ReturnJson JsonResult
	VersionRemainCounter := 0
	for key, _ := range DiffVersionResult {
		if key == 0 {
			DiffVersionResult[key].RemainCounter = 0
			VersionRemainCounter = DiffVersionResult[key].CreateCounter - DiffVersionResult[key].CloseCounter
		} else {
			DiffVersionResult[key].RemainCounter = VersionRemainCounter
			VersionRemainCounter = VersionRemainCounter + DiffVersionResult[key].CreateCounter - DiffVersionResult[key].CloseCounter
		}
		ReturnJson.Name = append(ReturnJson.Name, DiffVersionResult[key].Name)
		ReturnJson.CreateCounter = append(ReturnJson.CreateCounter, DiffVersionResult[key].CreateCounter)
		ReturnJson.CloseCounter = append(ReturnJson.CloseCounter, DiffVersionResult[key].CloseCounter)
		ReturnJson.RemainCounter = append(ReturnJson.RemainCounter, DiffVersionResult[key].RemainCounter)
		// log.Println(key, Value)
	}
	// log.Println(ReturnJson)
	return ReturnJson
}

func DiffVersionSoFarRemain(ProjectName string) JsonResult {
	ReturnJson := DiffVersion(ProjectName)
	soFarClose := 0
	for Key, _ := range ReturnJson.Name {
		tmp := ReturnJson.CloseCounter[Key]
		if Key == 0 {
			ReturnJson.CloseCounter[Key] = 0
		} else {
			ReturnJson.CloseCounter[Key] = soFarClose
		}
		soFarClose = soFarClose + tmp
	}
	return ReturnJson
}

func TYGHDiffVersionSoFarRemain(ProjectName string) (JsonResult, JsonResult) {
	ReturnJsonAPP, ReturnJsonWEB := TYGHDiffVersion(ProjectName)
	soFarAPPClose := 0
	soFarWEBClose := 0
	for Key, _ := range ReturnJsonAPP.Name {
		tmp := ReturnJsonAPP.CloseCounter[Key]
		if Key == 0 {
			ReturnJsonAPP.CloseCounter[Key] = 0
		} else {
			ReturnJsonAPP.CloseCounter[Key] = soFarAPPClose
		}
		soFarAPPClose = soFarAPPClose + tmp

		// soFarAPPClose = soFarAPPClose + ReturnJsonAPP.CloseCounter[Key]
		// ReturnJsonAPP.CloseCounter[Key] = soFarAPPClose
	}

	for Key, _ := range ReturnJsonWEB.Name {
		tmp := ReturnJsonWEB.CloseCounter[Key]
		if Key == 0 {
			ReturnJsonWEB.CloseCounter[Key] = 0
		} else {
			ReturnJsonWEB.CloseCounter[Key] = soFarWEBClose
		}
		soFarWEBClose = soFarWEBClose + tmp

		// soFarWEBClose = soFarWEBClose + ReturnJsonWEB.CloseCounter[Key]
		// ReturnJsonWEB.CloseCounter[Key] = soFarWEBClose
	}
	// log.Println(ReturnJsonAPP)
	// log.Println(ReturnJsonWEB)
	return ReturnJsonAPP, ReturnJsonWEB
}

func DiffDate(ProjectName string) JsonResult {
	var ReturnJson JsonResult
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	// Date Create

	DateRemainCreateSQL := "SELECT count(*),CreateDate FROM (SELECT *,DATE(CreatedTime) as CreateDate FROM Issues WHERE `Key` LIKE '" + ProjectName + "-%'  ORDER by CreatedTime) as b GROUP BY CreateDate"
	CreateDateRemains, err := db.Query(DateRemainCreateSQL)
	checkerr(DateRemainCreateSQL, err)
	var DateReaminResult []JiraIssuesDateRemain
	for CreateDateRemains.Next() {
		var (
			DateRemain    JiraIssuesDateRemain
			tmpCounter    int
			tmpCreateDate string
		)
		CreateDateRemains.Scan(&tmpCounter, &tmpCreateDate)
		DateRemain.Name = tmpCreateDate
		DateRemain.CreateCounter = tmpCounter
		DateReaminResult = append(DateReaminResult, DateRemain)
	}
	CreateDateRemains.Close()

	sortutil.AscByField(DateReaminResult, "Name")
	IssuesSQL := "SELECT Id FROM Issues WHERE `Status`='Closed' AND `Key` like '" + ProjectName + "-%' "
	Issuess, err := db.Query(IssuesSQL)
	checkerr(IssuesSQL, err)
	for Issuess.Next() {
		var tmpId string
		var tmpresult JiraIssuesDateRemain
		Issuess.Scan(&tmpId)
		HistorySQL := "SELECT date(Created) FROM `Historys` WHERE `Jid` = '" + tmpId + "' AND Tostring='Resolved' ORDER BY Sn DESC LIMIT 1"
		ResolveDate := QuerySingle(HistorySQL)
		appendflag := 0
		if ResolveDate == "" {
			HistorySQL := "SELECT date(Created) FROM `Historys` WHERE `Jid` = '" + tmpId + "' AND Tostring='Closed' ORDER BY Sn DESC LIMIT 1"
			ResolveDate = QuerySingle(HistorySQL)
		}
		// log.Println("Resolve date", ResolveDate, tmpId)
		for Key, Value := range DateReaminResult {
			if Value.Name == ResolveDate {
				appendflag++
				DateReaminResult[Key].CloseCounter++
			}
		}
		// log.Println(ResolveDate, tmpId, appendflag)
		if appendflag == 0 {

			tmpresult.Name = ResolveDate
			tmpresult.CloseCounter++
			DateReaminResult = append(DateReaminResult, tmpresult)
		}
	}
	Issuess.Close()
	db.Close()

	OldResultRemain := 0
	sortutil.AscByField(DateReaminResult, "Name")
	for key, Value := range DateReaminResult {
		if key == 0 {
			ReturnJson.RemainCounter = append(ReturnJson.RemainCounter, 0)
		} else {
			ReturnJson.RemainCounter = append(ReturnJson.RemainCounter, OldResultRemain)
		}
		OldResultRemain = (OldResultRemain + DateReaminResult[key].CreateCounter) - DateReaminResult[key].CloseCounter
		ReturnJson.Name = append(ReturnJson.Name, Value.Name)
		ReturnJson.CreateCounter = append(ReturnJson.CreateCounter, Value.CreateCounter)
		ReturnJson.CloseCounter = append(ReturnJson.CloseCounter, Value.CloseCounter)
	}
	// for key, value := range ReturnJson {
	// 	log.Println(ProjectName, key, value)
	// }
	log.Println(ReturnJson)
	return ReturnJson
}

func TYGHDiffDate(ProjectName string) (JsonResult, JsonResult) {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	// Date Create

	DateRemainCreateSQL := "SELECT DATE(CreatedTime) as CreateDate FROM Issues WHERE `Key` LIKE '" + ProjectName + "-%' group by CreateDate "
	CreateDateRemains, err := db.Query(DateRemainCreateSQL)
	checkerr(DateRemainCreateSQL, err)
	var APPDateReaminResult []JiraIssuesDateRemain
	var WEBDateReaminResult []JiraIssuesDateRemain
	for CreateDateRemains.Next() {
		var (
			DateRemain    JiraIssuesDateRemain
			tmpCreateDate string
		)
		CreateDateRemains.Scan(&tmpCreateDate)
		DateRemain.Name = tmpCreateDate
		APPDateReaminResult = append(APPDateReaminResult, DateRemain)
		WEBDateReaminResult = append(WEBDateReaminResult, DateRemain)
	}
	CreateDateRemains.Close()

	CreateIssueSQL := "SELECT Id,Version,DATE(CreatedTime) as createdate FROM `Issues` WHERE `Key` like '" + ProjectName + "-%' ORDER BY `Issues`.`Version` ASC"
	Issuess, err := db.Query(CreateIssueSQL)
	checkerr(CreateIssueSQL, err)
	for Issuess.Next() {
		var (
			tmpId    string
			tmpVerId string
			tmpDate  string
		)
		Issuess.Scan(&tmpId, &tmpVerId, &tmpDate)
		VersionsSQL := "SELECT Data FROM `Version` WHERE Id='" + tmpId + "' AND Enable='1'"
		VersionsROW, err := db.Query(VersionsSQL)
		checkerr(VersionsSQL, err)
		CreateVersionCounter := 0
		for VersionsROW.Next() {
			appappendflag := 0
			webappendflag := 0
			CreateVersionCounter++
			var versionID string
			var tmpresult JiraIssuesDateRemain
			VersionsROW.Scan(&versionID)
			VersionNameSQL := "SELECT Name FROM `Versions` WHERE Sn='" + versionID + "'"
			VersionName := QuerySingle(VersionNameSQL)
			// log.Println("Create", VersionName, tmpId, tmpDate)
			if strings.Contains(VersionName, "APP") {
				for Key, Value := range APPDateReaminResult {
					if Value.Name == tmpDate {
						appappendflag++
						APPDateReaminResult[Key].CreateCounter++
					}
				}
				if appappendflag == 0 {
					tmpresult.Name = tmpDate
					tmpresult.CreateCounter++
					APPDateReaminResult = append(APPDateReaminResult, tmpresult)
				}
			} else if strings.Contains(VersionName, "WEB") {
				for Key, Value := range WEBDateReaminResult {
					if Value.Name == tmpDate {
						webappendflag++
						WEBDateReaminResult[Key].CreateCounter++
					}
				}
				if webappendflag == 0 {
					log.Println("web", tmpDate)
					tmpresult.Name = tmpDate
					tmpresult.CreateCounter++
					WEBDateReaminResult = append(WEBDateReaminResult, tmpresult)
				}
			}
		}
		VersionsROW.Close()

		// log.Println(CreateVersionCounter)
		ResolvedIssueSQL := "SELECT Status as createdate FROM `Issues` WHERE `Key` like '" + ProjectName + "-%' AND Id='" + tmpId + "' ORDER BY `Issues`.`Version` ASC"
		Status := QuerySingle(ResolvedIssueSQL)
		if Status == "Closed" {
			HistorySQL := "SELECT date(Created) FROM `Historys` WHERE `Jid` = '" + tmpId + "' AND Tostring='Resolved' ORDER BY Sn DESC LIMIT 1"
			ResolveDate := QuerySingle(HistorySQL)
			if ResolveDate == "" {
				HistorySQL := "SELECT date(Created) FROM `Historys` WHERE `Jid` = '" + tmpId + "' AND Tostring='Closed' ORDER BY Sn DESC LIMIT 1"
				ResolveDate = QuerySingle(HistorySQL)
			}

			FixVersionsSQL := "SELECT Data FROM `Fixversion` WHERE Id='" + tmpId + "' AND Enable='1'"
			FixVersionsROW, err := db.Query(FixVersionsSQL)
			checkerr(FixVersionsSQL, err)
			for FixVersionsROW.Next() {
				var tmpresult JiraIssuesDateRemain
				appappendflag := 0
				webappendflag := 0
				var versionID string
				FixVersionsROW.Scan(&versionID)
				FixVersionSQL := "SELECT Name FROM `Fixversions` WHERE Sn='" + versionID + "'"
				FixVersionName := QuerySingle(FixVersionSQL)
				// log.Println("Close", FixVersionName, tmpId, ResolveDate, CreateVersionCounter)
				if strings.Contains(FixVersionName, "APP") && CreateVersionCounter >= 0 {
					for Key, Value := range APPDateReaminResult {
						if Value.Name == ResolveDate {
							appappendflag++
							APPDateReaminResult[Key].CloseCounter++
						}
					}
					if appappendflag == 0 {
						tmpresult.Name = ResolveDate
						tmpresult.CloseCounter++
						APPDateReaminResult = append(APPDateReaminResult, tmpresult)
					}
				} else if strings.Contains(FixVersionName, "WEB") && CreateVersionCounter >= 0 {
					for Key, Value := range WEBDateReaminResult {
						if Value.Name == ResolveDate {
							webappendflag++
							WEBDateReaminResult[Key].CloseCounter++
						}
					}
					if webappendflag == 0 {
						tmpresult.Name = ResolveDate
						tmpresult.CloseCounter++
						WEBDateReaminResult = append(WEBDateReaminResult, tmpresult)
					}
				}
				CreateVersionCounter--
			}
			FixVersionsROW.Close()
		}

	}
	Issuess.Close()
	db.Close()

	sortutil.AscByField(APPDateReaminResult, "Name")
	sortutil.AscByField(WEBDateReaminResult, "Name")

	WEBRemain := 0
	APPRemain := 0
	for key, _ := range WEBDateReaminResult {
		WEBRemain = (WEBRemain + WEBDateReaminResult[key].CreateCounter) - WEBDateReaminResult[key].CloseCounter
		// if key == 0 {
		// 	WEBDateReaminResult[key].RemainCounter = 0
		// } else {
		WEBDateReaminResult[key].RemainCounter = WEBRemain
		// }
		// WEBRemain = (WEBRemain + WEBDateReaminResult[key].CreateCounter) - WEBDateReaminResult[key].CloseCounter
	}
	for key, _ := range APPDateReaminResult {
		APPRemain = (APPRemain + APPDateReaminResult[key].CreateCounter) - APPDateReaminResult[key].CloseCounter
		// if key == 0 {
		// 	APPDateReaminResult[key].RemainCounter = 0
		// } else {
		APPDateReaminResult[key].RemainCounter = APPRemain
		// }
		// APPRemain = (APPRemain + APPDateReaminResult[key].CreateCounter) - APPDateReaminResult[key].CloseCounter
	}

	var ReturnJsonAPP JsonResult
	var ReturnJsonWEB JsonResult
	for key, _ := range WEBDateReaminResult {
		ReturnJsonWEB.Name = append(ReturnJsonWEB.Name, WEBDateReaminResult[key].Name)
		ReturnJsonWEB.CreateCounter = append(ReturnJsonWEB.CreateCounter, WEBDateReaminResult[key].CreateCounter)
		ReturnJsonWEB.CloseCounter = append(ReturnJsonWEB.CloseCounter, WEBDateReaminResult[key].CloseCounter)
		ReturnJsonWEB.RemainCounter = append(ReturnJsonWEB.RemainCounter, WEBDateReaminResult[key].RemainCounter)
	}
	for key, _ := range APPDateReaminResult {
		ReturnJsonAPP.Name = append(ReturnJsonAPP.Name, APPDateReaminResult[key].Name)
		ReturnJsonAPP.CreateCounter = append(ReturnJsonAPP.CreateCounter, APPDateReaminResult[key].CreateCounter)
		ReturnJsonAPP.CloseCounter = append(ReturnJsonAPP.CloseCounter, APPDateReaminResult[key].CloseCounter)
		ReturnJsonAPP.RemainCounter = append(ReturnJsonAPP.RemainCounter, APPDateReaminResult[key].RemainCounter)
	}

	// log.Println(ReturnJsonAPP, ReturnJsonWEB)
	return ReturnJsonAPP, ReturnJsonWEB
}
func PieChart(ProjectName string, TYPE string) JsonResultPie {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	var ReturnJson JsonResult
	var ReturnJsonPie JsonResultPie
	PrioritySQL := "Select `" + TYPE + "` FROM `Issues` WHERE `Key` like '" + ProjectName + "-%' group by `" + TYPE + "` "

	PriorityRows, err := db.Query(PrioritySQL)
	checkerr(PrioritySQL, err)
	for PriorityRows.Next() {
		var (
			tmpName string
		)
		PriorityRows.Scan(&tmpName)
		if tmpName != "" {
			ReturnJson.Name = append(ReturnJson.Name, tmpName)
			ReturnJsonPie.Name = append(ReturnJsonPie.Name, tmpName)
		}
	}
	PriorityRows.Close()
	totalcounter := 0
	for _, Value := range ReturnJson.Name {
		// log.Println(key, Value)
		tmpValue := Value
		TmpCounter := 0
		SQL := "Select count(*) FROM `Issues` WHERE `Key` like '" + ProjectName + "-%' AND `" + TYPE + "` = ?  "
		// log.Println(SQL, "="+tmpValue+"=")
		rows, err := db.Query(SQL, tmpValue)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&TmpCounter)
			if err != nil {
				log.Fatal(err)
			}
			// log.Println(TmpCounter)
		}
		rows.Close()
		ReturnJson.CreateCounter = append(ReturnJson.CreateCounter, TmpCounter)
		totalcounter = totalcounter + TmpCounter
	}

	for key, _ := range ReturnJson.Name {
		tmpremain := ((float64)(ReturnJson.CreateCounter[key]) * 100) / (float64)(totalcounter)
		// log.Println(tmpremain, ReturnJson.CreateCounter[key], totalcounter)
		ReturnJsonPie.RemainCounter = append(ReturnJsonPie.RemainCounter, tmpremain)
	}
	// log.Println(ReturnJsonPie)
	db.Close()
	// 	// log.Println(ReturnJson)
	return ReturnJsonPie
}

func IssueTimespent(ProjectName string) []JsonResultTable {
	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
	RemainIssueSQL := "SELECT a.`Key`, a.`Summary`, a.DiffDate from (SELECT *,DATEDIFF(DATE(`UpdatedTime`),DATE(CURDATE() )) AS DiffDate " +
		"FROM `Issues` WHERE `Key` LIKE '" + ProjectName + "-%' AND Status!='Closed' AND Resolution not like  '% Fix' ORDER BY UpdatedTime AND Status='Closed' )as a " +
		"ORDER by a.DiffDate"
	RemainIssues, err := db.Query(RemainIssueSQL)
	checkerr(RemainIssueSQL, err)
	var RemainIssueResult []JsonResultTable
	for RemainIssues.Next() {
		var tmpRemainIssue JsonResultTable
		RemainIssues.Scan(&tmpRemainIssue.Name, &tmpRemainIssue.Summary, &tmpRemainIssue.DiffDate)
		log.Println(tmpRemainIssue)
		RemainIssueResult = append(RemainIssueResult, tmpRemainIssue)

	}
	RemainIssues.Close()
	db.Close()
	return RemainIssueResult
}

// func TYGHDiffVersion(ProjectName string) (JsonResult, JsonResult) {
// 	var ReturnJsonAPP JsonResult
// 	var ReturnJsonWEB JsonResult
// 	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
// 	// ProjectName := "BABY"
// 	versSQL := "SELECT Id,Name from ( " +
// 		"SELECT f.Id,f.Name,f.Project FROM Fixversions as f LEFT join Versions as v on f.Id = v.Id UNION DISTINCT " +
// 		"SELECT v.Id,v.Name,v.Project FROM Fixversions as f RIGHT join Versions as v on f.Id = v.Id) as a where  Project='" + ProjectName + "'  group by a.Id " +
// 		"ORDER BY `a`.`Id` ASC"
// 	Vers, err := db.Query(versSQL)
// 	checkerr(versSQL, err)
// 	var DiffVersionResult []JiraIssuesStatusDiffVersion
// 	for Vers.Next() {
// 		var (
// 			IssueVersions JiraIssuesStatusDiffVersion
// 			tmpVerId      string
// 			tmpVerName    string
// 		)
// 		Vers.Scan(&tmpVerId, &tmpVerName)
// 		IssueVersions.Name = tmpVerName
// 		DiffVersionResult = append(DiffVersionResult, IssueVersions)
// 	}
// 	Vers.Close()

// 	sortutil.AscByField(DiffVersionResult, "Name")
// 	APPRemain := 0
// 	WEBRemain := 0

// 	for key, Value := range DiffVersionResult {

// 		VerSNSQL := " SELECT Sn FROM `Versions` WHERE `Project` = '" + ProjectName + "' AND Name ='" + Value.Name + "' "
// 		VerSN := QuerySingle(VerSNSQL)
// 		VerDataSQL := "SELECT Count(*) from Issues as i ,Versions as vs ,Version as v  where `Key` like '" + ProjectName + "-%' AND v.Data=vs.Sn AND vs.`Sn`='" + VerSN + "'  AND v.Id = i.Id AND v.`Enable`='1' AND i.Resolution not like '% Fix'"
// 		// log.Println(VerDataSQL)
// 		tmpCreateCounter := QuerySingle(VerDataSQL)
// 		Value.CreateCounter, _ = strconv.Atoi(tmpCreateCounter)

// 		FixVerSNSQL := " SELECT Sn FROM `Fixversions` WHERE `Project` = '" + ProjectName + "' AND Name ='" + Value.Name + "' "
// 		FixVerSN := QuerySingle(FixVerSNSQL)
// 		FixVerDataSQL := "SELECT count(*) from `Issues` as i inner join `Fixversions` as fs  inner join `Fixversion` as f where `Key` like '" + ProjectName + "-%' AND f.`Data`=fs.`Sn` AND fs.`Sn`='" + FixVerSN + "' AND i.`Status`='Closed' AND f.Id = i.Id AND f.`Enable`='1'  AND i.Resolution not like '% Fix'"
// 		// log.Println(FixVerDataSQL)
// 		tmpCloseCounter := QuerySingle(FixVerDataSQL)
// 		Value.CloseCounter, _ = strconv.Atoi(tmpCloseCounter)
// 		if strings.Contains(Value.Name, "WEB") {
// 			WEBRemain = (WEBRemain + Value.CreateCounter) - Value.CloseCounter
// 			Value.RemainCounter = WEBRemain
// 		} else if strings.Contains(Value.Name, "APP") {
// 			APPRemain = (APPRemain + Value.CreateCounter) - Value.CloseCounter
// 			Value.RemainCounter = APPRemain
// 		}
// 		log.Println(key, Value)
// 		if strings.Contains(Value.Name, "APP") {
// 			ReturnJsonAPP.Name = append(ReturnJsonAPP.Name, Value.Name)
// 			ReturnJsonAPP.CreateCounter = append(ReturnJsonAPP.CreateCounter, Value.CreateCounter)
// 			ReturnJsonAPP.CloseCounter = append(ReturnJsonAPP.CloseCounter, Value.CloseCounter)
// 			ReturnJsonAPP.RemainCounter = append(ReturnJsonAPP.RemainCounter, Value.RemainCounter)
// 		} else {
// 			ReturnJsonWEB.Name = append(ReturnJsonWEB.Name, Value.Name)
// 			ReturnJsonWEB.CreateCounter = append(ReturnJsonWEB.CreateCounter, Value.CreateCounter)
// 			ReturnJsonWEB.CloseCounter = append(ReturnJsonWEB.CloseCounter, Value.CloseCounter)
// 			ReturnJsonWEB.RemainCounter = append(ReturnJsonWEB.RemainCounter, Value.RemainCounter)
// 		}
// 	}

// 	db.Close()

// 	return ReturnJsonAPP, ReturnJsonWEB

// }

// func DiffVersion(ProjectName string) JsonResult {
// 	var ReturnJson JsonResult
// 	db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
// 	// ProjectName := "BABY"
// 	versSQL := "SELECT Id,Name from ( " +
// 		"SELECT f.Id,f.Name,f.Project FROM Fixversions as f LEFT join Versions as v on f.Id = v.Id UNION DISTINCT " +
// 		"SELECT v.Id,v.Name,v.Project FROM Fixversions as f RIGHT join Versions as v on f.Id = v.Id) as a where  Project='" + ProjectName + "'  group by a.Id " +
// 		"ORDER BY `a`.`Id` ASC"
// 	Vers, err := db.Query(versSQL)
// 	checkerr(versSQL, err)
// 	var DiffVersionResult []JiraIssuesStatusDiffVersion
// 	for Vers.Next() {
// 		var (
// 			IssueVersions JiraIssuesStatusDiffVersion
// 			tmpVerId      string
// 			tmpVerName    string
// 		)
// 		Vers.Scan(&tmpVerId, &tmpVerName)
// 		IssueVersions.Name = tmpVerName
// 		DiffVersionResult = append(DiffVersionResult, IssueVersions)
// 	}
// 	Vers.Close()

// 	sortutil.AscByField(DiffVersionResult, "Name")
// 	APPRemain := 0
// 	WEBRemain := 0
// 	ELSERemain := 0
// 	for key, Value := range DiffVersionResult {
// 		VerSNSQL := " SELECT Sn FROM `Versions` WHERE `Project` = '" + ProjectName + "' AND Name ='" + Value.Name + "' "
// 		VerSN := QuerySingle(VerSNSQL)
// 		VerDataSQL := "SELECT Count(*) from Issues as i ,Versions as vs ,Version as v  where `Key` like '" + ProjectName + "-%' AND v.Data=vs.Sn AND vs.`Sn`='" + VerSN + "'  AND v.Id = i.Id AND v.`Enable`='1' AND i.Resolution not like '% Fix'"
// 		// log.Println(VerDataSQL)
// 		tmpCreateCounter := QuerySingle(VerDataSQL)
// 		Value.CreateCounter, _ = strconv.Atoi(tmpCreateCounter)

// 		FixVerSNSQL := " SELECT Sn FROM `Fixversions` WHERE `Project` = '" + ProjectName + "' AND Name ='" + Value.Name + "' "
// 		FixVerSN := QuerySingle(FixVerSNSQL)
// 		FixVerDataSQL := "SELECT count(*) from `Issues` as i inner join `Fixversions` as fs  inner join `Fixversion` as f where `Key` like '" + ProjectName + "-%' AND f.`Data`=fs.`Sn` AND fs.`Sn`='" + FixVerSN + "' AND i.`Status`='Closed' AND f.Id = i.Id AND f.`Enable`='1'  AND i.Resolution not like '% Fix'"
// 		// log.Println(FixVerDataSQL)
// 		tmpCloseCounter := QuerySingle(FixVerDataSQL)
// 		Value.CloseCounter, _ = strconv.Atoi(tmpCloseCounter)
// 		if strings.Contains(Value.Name, "WEB") {
// 			WEBRemain = (WEBRemain + Value.CreateCounter) - Value.CloseCounter
// 			Value.RemainCounter = WEBRemain
// 		} else if strings.Contains(Value.Name, "APP") {
// 			APPRemain = (APPRemain + Value.CreateCounter) - Value.CloseCounter
// 			Value.RemainCounter = APPRemain
// 		} else {
// 			ELSERemain = (ELSERemain + Value.CreateCounter) - Value.CloseCounter
// 			Value.RemainCounter = ELSERemain
// 		}
// 		log.Println(key, Value)
// 		ReturnJson.Name = append(ReturnJson.Name, Value.Name)

// 		ReturnJson.CreateCounter = append(ReturnJson.CreateCounter, Value.CreateCounter)
// 		ReturnJson.CloseCounter = append(ReturnJson.CloseCounter, Value.CloseCounter)
// 		ReturnJson.RemainCounter = append(ReturnJson.RemainCounter, Value.RemainCounter)
// 	}

// 	db.Close()

// 	return ReturnJson

// }
// func WeekRemain(ProjectName string) JsonResult {
//Week remain
// var ReturnJson JsonResult
// db, err := sql.Open("mysql", "eli:eli@/Jira_Data")
// WeekRemainSQL := "SELECT a.week,a.createcounter ,b.closecounter from (" +
// 	"SELECT week , count(*) as createcounter FROM ( SELECT Sn,`Key`,CreatedTime , WEEK(CreatedTime)+1 as week FROM Issues " +
// 	" WHERE `Key` LIKE '" + ProjectName + "-%'  ORDER by CreatedTime ) as a  GROUP BY week ) as a  " +
// 	" INNER join (SELECT week,count(*) as closecounter  FROM ( SELECT Sn,`Key`,CreatedTime , WEEK(UpdatedTime)+1 as week FROM Issues " +
// 	" WHERE `Key` LIKE '" + ProjectName + "-%' AND Status='Closed'  ORDER by UpdatedTime ) as a  GROUP BY week  ) as b on a.week = b.week GROUP BY a.week "
// WeekRemains, err := db.Query(WeekRemainSQL)
// log.Println(WeekRemainSQL)
// checkerr(WeekRemainSQL, err)
// var WeekReaminResult []JiraIssuesDateRemain
// WeekRemainCounter := 0
// for WeekRemains.Next() {
// 	var (
// 		WeekRemain       JiraIssuesDateRemain
// 		tmpCreateCounter int
// 		tmpCloseCounter  int
// 		tmpName          string
// 	)
// 	WeekRemains.Scan(&tmpName, &tmpCreateCounter, &tmpCloseCounter)
// 	WeekRemain.Name = tmpName
// 	WeekRemain.CreateCounter = tmpCreateCounter
// 	WeekRemain.CloseCounter = tmpCloseCounter
// 	WeekRemainCounter = WeekRemainCounter + tmpCreateCounter - tmpCloseCounter
// 	WeekRemain.RemainCounter = WeekRemainCounter
// 	WeekReaminResult = append(WeekReaminResult, WeekRemain)
// }
// WeekRemains.Close()
// sortutil.AscByField(WeekReaminResult, "Name")
// for key, value := range WeekReaminResult {
// 	log.Println(key, value)
// 	ReturnJson.Name = append(ReturnJson.Name, value.Name)
// 	ReturnJson.CreateCounter = append(ReturnJson.CreateCounter, value.CreateCounter)
// 	ReturnJson.CloseCounter = append(ReturnJson.CloseCounter, value.CloseCounter)
// 	ReturnJson.RemainCounter = append(ReturnJson.RemainCounter, value.RemainCounter)
// }
// db.Close()
// return ReturnJson
// }
