// Welcome to the RazorFlow Dashbord Quickstart. Simply copy this "dashboard_quickstart"
// to somewhere in your computer/web-server to have a dashboard ready to use.
// This is a great way to get started with RazorFlow with minimal time in setup.
// However, once you're ready to go into deployment consult our documentation on tips for how to 
// maintain the most stable and secure 

StandaloneDashboard(function(db){
	
	db.setTabbedDashboardTitle ("HC300000 Dashboard");

	var db1 = new Dashboard();
	projectModel(db1,"TYGH"); //(Component,ProjectNmae)

	 var db2 = new Dashboard();
	 projectModel(db2,"BABY"); //(Component,ProjectNmae)

//	 var db3 = new Dashboard();
//	 projectModel(db3,"IOS"); //(Component,ProjectNmae)

	var db4 = new Dashboard();
	projectModel(db4,"MMHDRUG"); //(Component,ProjectNmae)

	db.addDashboardTab(db1, {
        title: 'TYGH project',
        // active: true
    });
    db.addDashboardTab(db2, {
        title: 'Baby project',
        active: true
    });

    db.addDashboardTab(db4, {
       title: 'MMH project',
    });

},{tabbed: true});


function projectModel(Component,ProjectNmae) {
	var mydata = JSON.parse(config);
    var kpi = new KPIComponent ();
    Component.addComponent(kpi);
    kpichart(kpi,"http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectNmae+"/LastUpdate");
	/** Show Pie chart  **/
	// var priority = new ChartComponent();
	// Component.addComponent (priority);
	// pieChar(priority,"Priority","http://"+location.hostname+":7108/"+ProjectNmae+"/Priority"); //(Component,CaptionNmae,Link)
	// /** Show Pie chart  **/
	// if (ProjectNmae != "IOS") {
	// 	var resolution = new ChartComponent();
	// 	Component.addComponent (resolution);
	// 	pieChar(resolution,"Resolution","http://"+location.hostname+":7108/"+ProjectNmae+"/Resolution"); //(Component,CaptionNmae,Link)
	// };

	// /** Show Pie chart  **/
	// var status = new ChartComponent();
	// Component.addComponent (status);
	// pieChar(status,"Status","http://"+location.hostname+":7108/"+ProjectNmae+"/Status"); //(Component,CaptionNmae,Link)

    var tableSummary = new TableComponent();
    Component.addComponent (tableSummary);
    ProjectSummaryTableChar(tableSummary,"Project summary ","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectNmae+"/ProjectSummary"); //(Component,CaptionNmae,Link)

	var duedatetable = new TableComponent();
	Component.addComponent (duedatetable);
    DueDateTableChar(duedatetable,"DueDate Remain Issue","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectNmae+"/DueDateRemain"); //(Component,CaptionNmae,Link)


    var table = new TableComponent();
    Component.addComponent (table);
    TableChar(table,"Remain Issue Time spent","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectNmae+"/IssueTimeSpent"); //(Component,CaptionNmae,Link)


/** Show Stacked Column chart  **/
	if(ProjectNmae == "TYGH"){
   
   		var so_far_status_APP = new ChartComponent();
    	so_far_status_APP.setDimensions (12, 4);
		Component.addComponent (so_far_status_APP);
		stackedColumnChar(so_far_status_APP,"JIRA issue status in diff. versions_APP","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectNmae+"/DiffVersionSoFarAPP");

		var so_far_status_WEB = new ChartComponent();
    	so_far_status_WEB.setDimensions (12, 4);
		Component.addComponent (so_far_status_WEB);
		stackedColumnChar(so_far_status_WEB,"JIRA issue status in diff. versions_WEB","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectNmae+"/DiffVersionSoFarWEB");

		/** Show Line chart  **/
	    var daily_status = new ChartComponent();
		Component.addComponent (daily_status);
		LineChar("TYGH",daily_status,"JIRA issue remain by time(Daily)","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectNmae+"/TYGHDiffDate");
    }else if (ProjectNmae == "MMHDRUG" ) {
  //   	var issue_status = new ChartComponent();
  //   	issue_status.setDimensions (6, 4);
		// Component.addComponent (issue_status);
		// stackedColumnChar(issue_status,"Issue_status","http://127.0.0.1:7108/"+ProjectNmae+"/DiffVersion");

		var so_far_status = new ChartComponent();
    	so_far_status.setDimensions (12, 4);
		Component.addComponent (so_far_status);
		stackedColumnChar(so_far_status,"JIRA issue status in diff. versions","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectNmae+"/DiffVersionSoFar");
    
		/** Show Line chart  **/
	    var daily_status = new ChartComponent();
		Component.addComponent (daily_status);
		LineChar("MMHDRUG",daily_status,"JIRA issue remain by time(Daily)","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectNmae+"/DiffDate");
    }else if (ProjectNmae == "BABY"  ) {
		var so_far_statusa = new ChartComponent();
    	so_far_statusa.setDimensions (12, 4);
		Component.addComponent (so_far_statusa);
		stackedColumnChar(so_far_statusa,"Android JIRA issue status in diff. versions","http://"+location.hostname+":"+mydata[0].port+"/api/BABY/DiffVersionSoFar");

		var so_far_status = new ChartComponent();
    	so_far_status.setDimensions (12, 4);
		Component.addComponent (so_far_status);
		stackedColumnChar(so_far_status,"Ios JIRA issue status in diff. versions","http://"+location.hostname+":"+mydata[0].port+"/api/IOS/DiffVersionSoFar");
    
		/** Show Line chart  **/
	    var daily_status = new ChartComponent();
		Component.addComponent (daily_status);
		LineChar("BABY",daily_status,"JIRA issue remain by time(Daily)","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectNmae+"/DiffDate");
    }

	// var weekly_status = new ChartComponent();
	// Component.addComponent (weekly_status);
	// LineChar(weekly_status,"Version_status(Weekly)","http://127.0.0.1:7108/"+ProjectNmae+"/WeekRemain");
}

function kpichart (Component,AddrLink) {
	// body...
    Component.setDimensions (3, 3);
    $.get(AddrLink, function (data) {
		Component.setValue ( data['Week'], {
		    numberPrefix : "Week "    
		});
	    Component.setCaption ( "Updated "+data['LastUpdate'] );
    });
}

function TableChar(Component,CaptionName,AddrLink) {
	var mydata = JSON.parse(config);
    Component.setCaption (CaptionName);
    Component.addColumn ("key", "Key" ,{columnWidth:"128",rawHTML:true});
    Component.addColumn ("summary", "Summary",{textAlign:"left",textBoldFlag:true});
    Component.addColumn ("priority", "Priority",{columnWidth:"128",rawHTML:true});
    Component.addColumn ("date", "Till today",{columnWidth:"128"});
    Component.addColumn ("assignee", "Assignee",{columnWidth:"128"});
    Component.setDimensions (12, 4);
    Component.lock ();

    var collectData = [];
    $.get(AddrLink, function (data) {
        // alert(data.length)
        for(i=0; i< data.length ;i++ )
        {
        	var parorityHTML = "" ;
        	if (data[i]['Priority'] == "Fatal" ) {
        		parorityHTML = "<font style='color:red;font-weight:900;font-size:18px' >" ;
        	}else if (data[i]['Priority'] == "Serious" ) {
        		parorityHTML = "<font style='color:orange;font-weight:900;font-size:18px' >" ;
        	}else if (data[i]['Priority'] == "Medium" ) {
        		parorityHTML = "<font  >" ;
        	}else if (data[i]['Priority'] == "Low" ) {
        		parorityHTML = "<font  style='color:green;font-weight:900' >" ;
        	};
        	parorityHTML = parorityHTML + data[i]['Priority'] + "</font>" ;

            Component.addRow ({
                "key": "<a href="+mydata[0].jirapath+data[i]['Key']+">"+data[i]['Key']+"</a>",
                "summary": data[i]['Summary'],
                "priority": parorityHTML,
                "date": Math.abs(data[i]['date']),
                "assignee": data[i]['assignee']
            });
        }
    });
    Component.unlock ();
}

function DueDateTableChar(Component,CaptionName,AddrLink) {
	var mydata = JSON.parse(config);
    Component.setCaption (CaptionName);
    Component.addColumn ("key", "Key" ,{columnWidth:"128",rawHTML:true});
    Component.addColumn ("summary", "Summary",{textAlign:"left",textBoldFlag:true});
    Component.addColumn ("priority", "Priority",{columnWidth:"128",rawHTML:true});
    Component.addColumn ("assignee", "Assignee",{columnWidth:"128"});
    Component.setDimensions (12, 2);
    Component.lock ();
	// Component.unlock ();

    var collectData = [];
    $.get(AddrLink, function (data) {
        // alert(data.length)
        if (data != null ) {
        	for(i=0; i< data.length ;i++ )
        	{
        		var parorityHTML = "" ;
        	if (data[i]['Priority'] == "Fatal" ) {
        		parorityHTML = "<font style='color:red;font-weight:900;font-size:18px' >" ;
        	}else if (data[i]['Priority'] == "Serious" ) {
        		parorityHTML = "<font style='color:orange;font-weight:900;font-size:18px' >" ;
        	}else if (data[i]['Priority'] == "Medium" ) {
        		parorityHTML = "<font  >" ;
        	}else if (data[i]['Priority'] == "Low" ) {
        		parorityHTML = "<font  style='color:green;font-weight:900' >" ;
        	};
        		parorityHTML = parorityHTML + data[i]['Priority'] + "</font>" ;

        		Component.addRow ({
        			"key": "<a href="+mydata[0].jirapath+data[i]['Key']+">"+data[i]['Key']+"</a>",
        			"summary": data[i]['Summary'],
        			"priority": parorityHTML,
        			"assignee": data[i]['assignee']
        		});
        	}        	
        };

    });
    // Component.lock ();
    Component.unlock ();
}


function ProjectSummaryTableChar(Component,CaptionName,AddrLink) {
	var mydata = JSON.parse(config);
    Component.setCaption (CaptionName);
    //`Project`,`Version`,`Scenario`,`DataCheck`,`Auto`,`BDI`,`Compatibility`,`Security`,`Other`,`Battery`,`Date`
    // Component.addColumn ("project", "Project" );
    Component.addColumn ("version", "Version" ,{rawHTML:true});
    Component.addColumn ("scenario", "Scenario",{rawHTML:true} );
    Component.addColumn ("dataCheck", "DataCheck",{rawHTML:true} );
    Component.addColumn ("auto", "Auto",{rawHTML:true} );
    Component.addColumn ("bDI", "BDI",{rawHTML:true} );
    Component.addColumn ("compatibility", "Compatibility",{rawHTML:true} );
    Component.addColumn ("other", "Other",{rawHTML:true} );
    Component.addColumn ("battery", "Battery",{rawHTML:true} );
    Component.addColumn ("date", "Date",{rawHTML:true} );
    Component.setDimensions (9, 3);
    Component.lock ();

    $.get(AddrLink, function (data) {
        // alert(data.length)
        for(i=0; i< data.length ;i++ )
        {
            Component.addRow ({
                // "project": data[i]['Project'],
                "version": "<a href="+data[i]['Version']+">"+data[i]['Version']+"</a>",
                // "scenario": data[i]['Scenario'],
                "scenario": TextColor("Scenario", data[i]['Scenario'] ),
                "dataCheck":  TextColor("DataCheck", data[i]['DataCheck'] ),
                // "dataCheck":  data[i]['DataCheck'],
                "auto": TextColor("Auto", data[i]['Auto'] ),
                "bDI": TextColor("BDI", data[i]['BDI'] ),
                "compatibility": TextColor("Compatibility", data[i]['Compatibility'] ),
                "other": TextColor("Other", data[i]['Other'] ),
                "battery": TextColor("Battery", data[i]['Battery'] ),
                "date": data[i]['Date']
            });
        }
    });
    Component.unlock ();
}

function pieChar(Component,CaptionName,AddrLink) {
    
    Component.setCaption(CaptionName);
	Component.setDimensions (3, 3);
	Component.lock ();	
	
	
	var collectData = [];	
	$.get(AddrLink, function (data) {
		Component.setLabels (data['Name']); // You can also use data.categories
        Component.addSeries (CaptionName, "items", data['RemainCounter'],{
        	seriesDisplayType: 'pie',
        	numberFormatFlag: true, 
        	numberDecimalPoints: 2,
        	//dataType: "number",
        	//numberHumanize: false, 
        	//numberForceDecimals: false,
        	//seriesStacked : true,
        	//seriesHiddenFlag: true
        	numberSuffix: "%"
        });
        // Don't forget to call unlock or the data won't be displayed
        Component.unlock ();
	});
}

function stackedColumnChar(Component,CaptionName,AddrLink) {
    
    Component.setCaption(CaptionName);
	Component.lock ();	
	
	$.get(AddrLink, function (data) {
	// $.get("http://127.0.0.1:7108/lookup/BABY", function (data) {
		Component.setLabels (data['Name']); // You can also use data.categories
        Component.addSeries ("Create", "Create", data['CreateCounter'],{
        	seriesStacked: true
        });
        Component.addSeries ("Close", "Close", data['CloseCounter'],{
        	seriesStacked: true
        });

        Component.addSeries ("Remain", "Remain", data['RemainCounter'],{
        	seriesStacked: true
        });

        // Don't forget to call unlock or the data won't be displayed
        Component.unlock ();
	});
}

function LineChar(ProjectNmae,Component,CaptionName,AddrLink) {
    
    Component.setCaption(CaptionName);
	Component.setDimensions (12, 4);
	
	if(ProjectNmae == "TYGH"){
		Component.lock ();	
		$.get(AddrLink, function (data) {
			
			Component.setLabels (data[0]['Name']); // You can also use data.categories

	        Component.addSeries ("APP", "APP", data[0]['RemainCounter'],{
	        	seriesDisplayType: "line"

	        });
	        Component.addSeries ("WEB", "WEB", data[1]['RemainCounter'],{
	        	seriesDisplayType: "line"

	        });
	        // Don't forget to call unlock or the data won't be displayed
		});
        Component.unlock ();
	}else if (ProjectNmae == "BABY"){
		Component.lock ();	
		$.get(AddrLink, function (data) {
			
			Component.setLabels (data[0]['Name']); // You can also use data.categories

	        Component.addSeries ("Android", "Android", data[0]['RemainCounter'],{
	        	seriesDisplayType: "line"

	        });
	        Component.addSeries ("IOS", "IOS", data[1]['RemainCounter'],{
	        	seriesDisplayType: "line"

	        });
	        // Don't forget to call unlock or the data won't be displayed
	        // Component.unlock ();
		});
        Component.unlock ();
	}else if (ProjectNmae == "MMHDRUG"){
		Component.lock ();	
			$.get(AddrLink, function (data) {
			
			Component.setLabels (data['Name']); // You can also use data.categories
	        Component.addSeries ("APP", "APP", data['RemainCounter'],{
	        	seriesDisplayType: "line"
	        });
	       
	        // Don't forget to call unlock or the data won't be displayed
		});
        Component.unlock ();
	}
    // Component.unlock ();
}

function TextColor(ColumnName,ColumnValue){
	var parorityHTML = "" ;
	var passrate = 80 ;
	switch(ColumnName) {
		case "Scenario":
		passrate = 98 ;
		break;
		case "DataCheck":
		passrate = 100 ;
		break;
		case "Auto":
		passrate = 90 ;
		break;
		case "BDI":
		passrate = 90 ;
		break;
		case "Compatibility":
		passrate = 95 ;
		break;
		case "Other":
		passrate = 80 ;
		break;
		case "Battery":
		passrate = 80 ;
		break;

		default:
		passrate = 100 ;
	}
	console.log(passrate +" "+ColumnValue+" "+ColumnName);
	if ( ColumnValue < passrate ) {
		parorityHTML = "<font style='color:red;font-weight:900;font-size:18px' >" ;
	}else  {
		parorityHTML = "<font  >" ;
	};
	parorityHTML = parorityHTML + ColumnValue + "</font>" ;
	return parorityHTML;
}







