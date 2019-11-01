# ===========================================================================
# =============== Clase que genera los reportes necesarios ==================
# ===========================================================================
#
#################################################################
# Software created by: Josue Pradas Sacristan for Sogeti
# Contact: josue.pradas-sacristan@sogeti.com
#################################################################
#
import maya
from datetime import date
from docxtpl import DocxTemplate, RichText
import xlsxwriter
import jinja2
import unirest
import resources
import os.path

def generate_detailed_report(project, outputname, sqc) :
    doc = None

    print "> Calculating language metrics..."
    ncloc = sqc.get_loc(project[1]["key"])
    languages = sqc.get_loc_lang_metrics(project[1]["key"], ncloc)
    ncloc = resources.span_number_format(int(ncloc))

    print "> Retrieving issues..."
    vulnerabilities_num = sqc.get_issues_on_type("VULNERABILITY,SECURITY_HOTSPOT", project[1]["key"])
    vulnerabilities = resources.span_number_format(int(vulnerabilities_num))
    bugs_num = sqc.get_issues_on_type("BUG", project[1]["key"])
    bugs = resources.span_number_format(int(bugs_num))
    code_smells_num = sqc.get_issues_on_type("CODE_SMELL", project[1]["key"])
    code_smells = resources.span_number_format(int(code_smells_num))

    print "> Retrieving components..."
    duplicated_blocks = resources.span_number_format(int(sqc.get_duplicated_blocks(project[1]["key"])))
    duplicity = resources.span_number_format(float(sqc.get_duplicated_density(project[1]["key"])))

    print "> Retrieving violated rules..."
    print "        -> BUGS"
    bug_violations = sqc.get_violations(languages, "BUG", project[1]["key"])
    print "        -> VULNERABILITIES"
    vuln_violations = sqc.get_violations(languages, "VULNERABILITY,SECURITY_HOTSPOT", project[1]["key"])
    print "        -> CODE SMELL"
    cs_violations = sqc.get_violations(languages, "CODE_SMELL", project[1]["key"])

    report_date_str = project[1]["lastAnalysisDate"]
    dt = maya.parse(report_date_str).datetime()
    report_date = str(dt.date().day) + "-" + resources.meses[str(dt.date().month)] + "-" + str(dt.date().year) + " " + str(dt.time())
    hoy = date.today()
    today = str(hoy.day) + " de " + resources.meses[str(hoy.month) + "c"] + " de " + str(hoy.year)


    try:
        doc = DocxTemplate("templates/static_detail_report_template.docx")
    except:
        print "> Error. .docx template file not found."
    jinja_env = jinja2.Environment(extensions=['jinja2.ext.loopcontrols'])
    if doc:
        rt = RichText()
        rt.add(project[0], url_id=doc.build_url_id('http://srv-analiza/sonar/dashboard?id=' + project[1]["key"]))

        context = {
            'today': today,
            'project_name': project[0],
            'rt': rt,
            'report_date': report_date,
            'ncloc': ncloc,
            'lang_metrics': languages,
            'vulnerabilities': vulnerabilities,
            'bugs': bugs,
            'code_smells': code_smells,
            'duplicated_blocks': duplicated_blocks,
            'duplicity': duplicity,
            'bugs_violations': bug_violations,
            'vuln_violations': vuln_violations,
            'cs_violations': cs_violations,
        }

        doc.render(context, jinja_env)
        doc.save(outputname + ".docx")

    generate_charts(project, languages, sqc, vulnerabilities_num, bugs_num, code_smells_num, outputname)



def generate_executive_report(project, outputname, sqc) :
    doc = None

    print "> Calculating language metrics..."
    ncloc = sqc.get_loc(project[1]["key"])
    languages = sqc.get_loc_lang_metrics(project[1]["key"], ncloc)

    print "> Retrieving issues..."
    vulnerabilities = sqc.get_issues_on_type("VULNERABILITY,SECURITY_HOTSPOT", project[1]["key"])
    bugs = sqc.get_issues_on_type("BUG", project[1]["key"])
    code_smells = sqc.get_issues_on_type("CODE_SMELL", project[1]["key"])

    context = {
        'project_name': project[0],
        'report_date': project[1]["lastAnalysisDate"],
        'ncloc': ncloc,
        'lang_metrics': languages,
    }

    try:
        doc = DocxTemplate("templates/static_executive_report_template.docx")
    except:
        print "> Error. .docx template file not found."
    jinja_env = jinja2.Environment(extensions=['jinja2.ext.loopcontrols'])
    if doc:
        doc.render(context, jinja_env)
        doc.save(outputname + ".docx")

    if(os.path.isfile("./" + outputname + ".xlsx")) :
        print ".xlsx file already exists"
    else :
        print "> Creating charts file..."
        generate_charts(project, languages, sqc, vulnerabilities, bugs, code_smells, outputname)


def generate_charts(project, languages, sqc, vulnerabilities, bugs, code_smells, outputname):
    print "> Retrieving chart data..."
    workbook = xlsxwriter.Workbook(outputname + '.xlsx')
    worksheet = workbook.add_worksheet()
    bold = workbook.add_format({'bold': 1})

# ===============================================
# COLUMN DATATABLE CHART

    data = [
        ['Language','Vulnerability', 'Bug', 'Code Smell'],
    ]
    # print languages #TODO cambiar el nombre de los lenguajes despues de las llamadas
    for language in languages :
        row = [language[0], sqc.get_issues_on_type_and_language("VULNERABILITY,SECURITY_HOTSPOT", project[1]["key"], language[3]), \
        sqc.get_issues_on_type_and_language("BUG", project[1]["key"], language[3]), \
        sqc.get_issues_on_type_and_language("CODE_SMELL", project[1]["key"], language[3])]
        data.append(row)

    i = 1
    for row in data:
        worksheet.write_row('A' + str(i), row)
        i += 1

    bars = workbook.add_chart({'type': 'column'})

    j = 2
    for language in languages:
        bars.add_series({
            'name': '=Sheet1!$A$'+str(j),
            'categories': '=Sheet1!$B$1:$D$1',
            'values': '=Sheet1!$B$'+ str(j) + ':$D$' + str(j)
        })
        j += 1


    bars.set_title({'name': 'Problemas por tipo y lenguaje'})
    bars.set_table({'show_keys': True})
    bars.set_legend({'position': 'none'})

    worksheet.insert_chart('G2', bars, {'x_offset': 25, 'y_offset': 10})
# ====================================================
# PIE CHART

    data = [
        ['Issues', 'Number of Issues'],
        ['Vulnerabilities', vulnerabilities],
        ['Bugs', bugs],
        ['Code Smells', code_smells]
    ]

    i += 1
    j = i + 1
    for row in data :
        worksheet.write_row('A' + str(i), row)
        i += 1

    pie = workbook.add_chart({'type': 'pie'})

    pie.add_series({
        'name': ' ',
        'categories': '=Sheet1!$A$' + str(j) + ':$A$' + str(j+2),
        'values': '=Sheet1!$B$' + str(j) + ':$B$' + str(j+2),
        'data_labels': {'percentage':True, 'position':'outside_end'}
    })

    pie.set_style(10)
    worksheet.insert_chart('P2', pie, {'x_offset': 25, 'y_offset': 10})

# ======================================================
# STACKED COLUMN CHART
    data = [
        ['Lenguaje', 'Cumplidas', 'Incumplidas']
    ]

    for language in languages :
        #rules = sqc.get_num_rules_by_qprofile(resources.Qprofiles[language[0]])
        rules = sqc.get_num_rules_by_qprofile(sqc.get_qprofile_key_by_lang(project[1]["key"], language[3]))
        violations = sqc.get_num_violations_by_lang(project[1]["key"], language[3])
        row = [language[0], int(rules - violations) , violations]
        data.append(row)

    i += 1
    j = i + 1
    for row in data :
        worksheet.write_row('A' + str(i), row)
        i += 1

    stacked = workbook.add_chart({'type': 'column', 'subtype': 'percent_stacked'})

    for x in ('B','C'):
        stacked.add_series({
            'name': '=Sheet1!$'+ x +'$'+ str(j-1),
            'categories': '=Sheet1!$A$'+ str(j) +':$A$' + str(j + len(languages) - 1),
            'values': '=Sheet1!$'+ x +'$'+ str(j) +':$'+ x +'$' + str(j + len(languages) - 1)
        })

    stacked.set_title({'name': 'Nivel de cumplimiento - Nivel 1'})
    stacked.set_style(13)

    worksheet.insert_chart('G20', stacked, {'x_offset': 25, 'y_offset': 10})

# ======================================================
# TRENDING HISTORY CHART

    history = sqc.get_project_analyses_issues(project[1]["key"])
    duplicity = sqc.get_project_analyses_duplicity(project[1]["key"])
    data = []
    dup = []
    dates_row = list()
    dates_row.append("Types")
    long = 0
    for dates in history[0]["history"]:
        dt = maya.parse(dates["date"]).datetime()
        report_date = str(dt.date().day) + "/" + resources.meses[str(dt.date().month)] + "/" + str(dt.date().year)
        dates_row.append(report_date)
        long += 1

    data.append(dates_row)
    dup.append(dates_row)
    for measure in history:
        dates_row = list()
        dates_row.append(measure["metric"])
        for value in measure["history"]:
            dates_row.append(float(value["value"]))
        data.append(dates_row)

    for statics in duplicity:
        dates_row = list()
        dates_row.append(statics["metric"])
        for value in statics["history"]:
            dates_row.append(float(value["value"]))
        dup.append(dates_row)

    i += 1
    j = i + 1
    # print data
    for row in data:
        worksheet.write_row('A' + str(i), row)
        i += 1

    bars_history = workbook.add_chart({'type': 'column'})

    cont = 1
    for x in ('B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'):
        if cont > long:
            break
        bars_history.add_series({
            'name': '=Sheet1!$'+ x +'$'+ str(j-1),
            'categories': '=Sheet1!$A$'+ str(j) +':$A$' + str(j+2),
            'values': '=Sheet1!$'+ x +'$'+ str(j) +':$'+ x +'$' + str(j + 2),
            'data_labels': {'value': True},
        })
        cont += 1

    bars_history.set_title({'name': 'Tendencia historica por tipo'})
    worksheet.insert_chart('P20', bars_history, {'x_offset': 25, 'y_offset': 10})

    i += 1
    j = i + 1
    h = 0
    # print dup
    columns = ['A','B','C','D']
    #,'E','F','G','H','I','J','K','L','M','N','O','P','Q','R','S','T','U','V','W','X','Y','Z']
    for column in dup:
        worksheet.write_column(columns[h] + str(j-1), column)
        h += 1

    i += long
    duplicated_blocks = workbook.add_chart({'type': 'line'})
    duplicated_lines_density = workbook.add_chart({'type': 'line'})

    duplicated_blocks.add_series({
        'name': '=Sheet1!$B$'+ str(j-1),
        'categories': '=Sheet1!$A$'+ str(j) +':$A$' + str(j + 1 ),
        'values': '=Sheet1!$B$'+ str(j) +':$B$' + str(j + 1)
    })

    duplicated_lines_density.add_series({
        'name': '=Sheet1!$C$'+ str(j-1),
        'categories': '=Sheet1!$A$'+ str(j) +':$A$' + str(j + 1 ),
        'values': '=Sheet1!$C$'+ str(j) +':$C$' + str(j + 1)
    })

    duplicated_blocks.set_title({'name': 'Bloques duplicados'})
    worksheet.insert_chart('G37', duplicated_blocks, {'x_offset': 25, 'y_offset': 10})
    duplicated_lines_density.set_title({'name': 'Duplicidad'})
    worksheet.insert_chart('P37', duplicated_lines_density, {'x_offset': 25, 'y_offset': 10})

    workbook.close()
