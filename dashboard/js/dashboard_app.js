// Welcome to the RazorFlow Dashbord Quickstart. Simply copy this "dashboard_quickstart"
// to somewhere in your computer/web-server to have a dashboard ready to use.
// This is a great way to get started with RazorFlow with minimal time in setup.
// However, once you're ready to go into deployment consult our documentation on tips for how to 
// maintain the most stable and secure 

StandaloneDashboard(function(db){
	
	db.setTabbedDashboardTitle ("HC300000 Dashboard");

	var db1 = new Dashboard();
	projectModel(db1,"TYGH"); //(Component,ProjectName)

	 var db2 = new Dashboard();
	 projectModel(db2,"BABY"); //(Component,ProjectName)

	var db4 = new Dashboard();
	projectModel(db4,"MMHDRUG"); //(Component,ProjectName)


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
        // active: true
    });


},{tabbed: true});

function sleep(delay) {
    var start = new Date().getTime();
    while (new Date().getTime() < start + delay);
}

function projectModel(Component,ProjectName) {
	var mydata = JSON.parse(config);
    var kpi = new KPIComponent ();
    Component.addComponent(kpi);
    kpichart(kpi,"http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectName+"/LastUpdate");

    var tableSummary = new TableComponent();
    Component.addComponent (tableSummary);
    ProjectSummaryTableChar(tableSummary,"Project summary ","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectName+"/ProjectSummary",ProjectName); //(Component,CaptionNmae,Link)

	var duedatetable = new TableComponent();
	Component.addComponent (duedatetable);
    DueDateTableChar(duedatetable,"DueDate Remain Issue","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectName+"/DueDateRemain"); //(Component,CaptionNmae,Link)


    var table = new TableComponent();
    Component.addComponent (table);
    TableChar(table,"Remain Issue Time spent","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectName+"/IssueTimeSpent"); //(Component,CaptionNmae,Link)


    /** Show Stacked Column chart  **/
	if(ProjectName == "TYGH"){
   
   		var so_far_status_APP = new ChartComponent();
    	so_far_status_APP.setDimensions (12, 4);
		Component.addComponent (so_far_status_APP);
		stackedColumnChar(so_far_status_APP,"JIRA issue status in diff. versions_APP","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectName+"/DiffVersionSoFarAPP");

		var so_far_status_WEB = new ChartComponent();
    	so_far_status_WEB.setDimensions (12, 4);
		Component.addComponent (so_far_status_WEB);
		stackedColumnChar(so_far_status_WEB,"JIRA issue status in diff. versions_WEB","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectName+"/DiffVersionSoFarWEB");

		/** Show Line chart  **/
        /** show week line chart **/
	    var daily_status = new ChartComponent();
		Component.addComponent (daily_status);
		ColumnChar("TYGH",daily_status,"JIRA issue remain by time(Weekly)","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectName+"/DiffWeek");

    }else if (ProjectName == "MMHDRUG" ) {

		var so_far_status = new ChartComponent();
    	so_far_status.setDimensions (12, 4);
		Component.addComponent (so_far_status);
		stackedColumnChar(so_far_status,"JIRA issue status in diff. versions","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectName+"/DiffVersionSoFar");
    
		// /** Show Line chart  **/
	    var daily_status = new ChartComponent();
		Component.addComponent (daily_status);
		ColumnChar("MMHDRUG",daily_status,"JIRA issue remain by time(Daily)","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectName+"/DiffWeek");

    }else if (ProjectName == "BABY"  ) {
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
        daily_status.setCaption("JIRA issue remain by time(Daily)");
		ColumnChar("BABY",daily_status,"JIRA issue remain by time(Daily)","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectName+"/DiffWeek");

    }
    var priority = new ChartComponent();
    Component.addComponent (priority);
    pieChar(priority,"Priority","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectName+"/Priority"); //(Component,CaptionNmae,Link)
    /** Show Pie chart  **/
    var status = new ChartComponent();
    Component.addComponent (status);
    pieChar(status,"Status","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectName+"/Status"); //(Component,CaptionNmae,Link)
    /** Show Pie chart  **/
    var resolution = new ChartComponent();
    Component.addComponent (resolution);
    pieChar(resolution,"Resolution","http://"+location.hostname+":"+mydata[0].port+"/api/"+ProjectName+"/Resolution"); //(Component,CaptionNmae,Link)

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
    Component.setDimensions (12, 4);
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


function ProjectSummaryTableChar(Component,CaptionName,AddrLink,ProjectName) {
	var mydata = JSON.parse(config);
    Component.setCaption (CaptionName);
    //`Project`,`Version`,`Scenario`,`DataCheck`,`Auto`,`BDI`,`Compatibility`,`Security`,`Other`,`Battery`,`Date`
    // Component.addColumn ("project", "Project" );
    Component.setRowsPerPage(6);
    Component.addColumn ("version", "Version" ,{rawHTML:true});
    Component.addColumn ("scenario", "Scenario",{rawHTML:true,columnWidth:"80"} );
    Component.addColumn ("auto", "Auto",{rawHTML:true,columnWidth:"80"} );
    Component.addColumn ("bDI", "BDI",{rawHTML:true,columnWidth:"80"} );
    Component.addColumn ("compatibility", "Compatibility",{rawHTML:true,columnWidth:"80"} );
    Component.addColumn ("other", "Other",{rawHTML:true,columnWidth:"80"} );
    Component.addColumn ("battery", "Battery",{rawHTML:true,columnWidth:"80"} );
    Component.addColumn ("date", "Date",{rawHTML:true,columnWidth:"128"} );
    Component.setDimensions (9, 3);
    Component.lock ();

    $.get(AddrLink, function (data) {
        // alert(data.length)
        for(i=0; i< data.length ;i++ )
        {
            Component.addRow ({
                // "project": data[i]['Project'],
                "version": "<a href="+"upload/"+ProjectName+"/"+data[i]['FileName']+">"+data[i]['Version']+"</a>",
                // "scenario": data[i]['Scenario'],
                "scenario": TextColor("Scenario", data[i]['Scenario'] ),
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

function ColumnChar(ProjectName,Component,CaptionName,AddrLink) {
    
    Component.setCaption(CaptionName);
	Component.setDimensions (12, 4);
	
	if(ProjectName == "TYGH"){
        Component.setCaption("JIRA issue remain by time(Weekly)");
        Component.setDimensions (12, 4);
        Component.lock ();
         $.get(AddrLink, function (data) {
        Component.setLabels (data[0]['Name']); // You can also use data.categories
        Component.addSeries ("APP", "APP", data[0]['RemainCounter'],{seriesDisplayType: "column" });
        Component.addSeries ("WEB", "WEB", data[1]['RemainCounter'],{seriesDisplayType: "column" });
        Component.unlock ();
        });
        Component.addDrillStep (function (done, params, updatedComponent) {
            updatedComponent.lock();
            Component.setCaption("JIRA issue remain by time(Daily)");
            $.get("http://10.116.136.13:7108/api/"+ProjectName+"/WeekToDate/"+params.label, function (data) {
                updatedComponent.setLabels (data[0]['Name']); // You can also use data.categories
                updatedComponent.addSeries ("APP", "APP", data[0]['RemainCounter'],{ seriesDisplayType: "column" });
                updatedComponent.addSeries ("WEB", "WEB", data[1]['RemainCounter'],{ seriesDisplayType: "column" });
                done();
                updatedComponent.unlock ();
            });
        });
	}else if (ProjectName == "BABY"){
        Component.setCaption("JIRA issue remain by time(Weekly)");
        Component.setDimensions (12, 4);
        Component.lock ();
         $.get(AddrLink, function (data) {
        Component.setLabels (data[0]['Name']); // You can also use data.categories
        Component.addSeries ("APP", "APP", data[0]['RemainCounter'],{seriesDisplayType: "column" });
        Component.addSeries ("WEB", "WEB", data[1]['RemainCounter'],{seriesDisplayType: "column" });
        Component.unlock ();
        });
        Component.addDrillStep (function (done, params, updatedComponent) {
            updatedComponent.lock();
            Component.setCaption("JIRA issue remain by time(Daily)");
            $.get("http://10.116.136.13:7108/api/"+ProjectName+"/WeekToDate/"+params.label, function (data) {
                updatedComponent.setLabels (data[0]['Name']); // You can also use data.categories
                updatedComponent.addSeries ("Android", "Android", data[0]['RemainCounter'],{ seriesDisplayType: "column" });
                updatedComponent.addSeries ("IOS", "IOS", data[1]['RemainCounter'],{ seriesDisplayType: "column" });
                done();
                updatedComponent.unlock ();
            });
        });
	}else if (ProjectName == "MMHDRUG"){
        Component.setCaption("JIRA issue remain by time(Weekly)");
        Component.setDimensions (12, 4);
        Component.lock ();
         $.get(AddrLink, function (data) {
        Component.setLabels (data['Name']); // You can also use data.categories
        Component.addSeries ("APP", "APP", data['RemainCounter'],{seriesDisplayType: "column" });
        Component.unlock ();
        });
        Component.addDrillStep (function (done, params, updatedComponent) {
            updatedComponent.lock();
            Component.setCaption("JIRA issue remain by time(Daily)");
            $.get("http://10.116.136.13:7108/api/"+ProjectName+"/WeekToDate/"+params.label, function (data) {
                updatedComponent.setLabels (data['Name']); // You can also use data.categories
                updatedComponent.addSeries ("APP", "APP", data['RemainCounter'],{ seriesDisplayType: "column" });
                done();
                updatedComponent.unlock ();
            });
        });
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
	// console.log(passrate +" "+ColumnValue+" "+ColumnName);
    if ( ColumnValue == "0" ) {
        parorityHTML = "<font style='color:black;font-weight:100;font-size:6px' >" ;
        parorityHTML = parorityHTML + "NA" + "</font>" ;
    }else {
        if ( ColumnValue < passrate ) {
            parorityHTML = "<font style='color:red;font-weight:900;font-size:18px' >" ;
        }else  {
            parorityHTML = "<font  >" ;
        };        
        parorityHTML = parorityHTML + ColumnValue + "</font>" ;
    };

	return parorityHTML;
}







