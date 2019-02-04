#=================================================================
#======== Clase que controla el acceso al API Sonarqube ==========
#=================================================================
#
#################################################################
# Software created by: Josue Pradas Sacristan for Sogeti
# Contact: josue.pradas-sacristan@sogeti.com
#################################################################
#

import unirest
import sys
import resources

class SonarQubeAPI :

    def __init__(self) :
        self.user = resources.user
        self.psw = resources.pswd
        self.url = resources.url_base
        self.headers = resources.headers
        self.test_connection()

    def test_connection(self) :
        try :
            response = unirest.get(self.url + "api/authentication/validate", headers=self.headers, auth=(self.user, self.psw))
        except Exception as e :
            print "Server unreacheable. Check provided URL: " + self.url
            sys.exit()

    def get_organizations(self) :
        response = unirest.get(self.url + "api/organizations/search", headers=self.headers, auth=(self.user, self.psw))
        o = list()
        try :
            organizations = response.body["organizations"]
            for org in organizations:
                o.append((org["name"], org["key"]))
            return o
        except :
            return o

    def get_projects(self, org) :
        response = unirest.get(self.url + "api/projects/search", headers=self.headers, auth=(self.user, self.psw), params={"organization": org})
        p = list()
        try :
            projects = response.body["components"]
            for project in projects:
                p.append((project["name"], project))
            return p
        except :
            return p

    def get_loc(self, project) :
        response = unirest.get(self.url + "api/measures/component", headers=self.headers, auth=(self.user, self.psw), params={"component": project, "metricKeys": "ncloc"})
        return response.body["component"]["measures"][0]["value"]

    def get_loc_lang_metrics(self, project, ncloc) :
        response = unirest.get(self.url + "api/measures/component", headers=self.headers, auth=(self.user, self.psw), params={"component": project, "metricKeys": "ncloc_language_distribution"})
        languages = response.body["component"]["measures"][0]["value"].split(";")
        m = list()
        for language in languages :
            s = language.split("=")
            dividendo = float(s[1])
            divisor = float(ncloc)
            percent = round(float(dividendo/divisor)*100, 2)
            m.append((resources.languages[s[0]], s[1], resources.span_number_format(float(percent)), s[0], resources.span_number_format(int(s[1]))))
        return m

    def get_issues_number(self, project) :
        response = unirest.get(self.url + "api/issues/search", headers=self.headers, auth=(self.user, self.psw), params={"projects": project, "statuses": "OPEN,REOPENED,CONFIRMED", "types": "BUG,VULNERABILITY,CODE_SMELL,SECURITY_HOTSPOT"})
        return response.body["total"]

    def get_issues_on_type(self, type, project) :
        response = unirest.get(self.url + "api/issues/search", headers=self.headers, auth=(self.user, self.psw), params={"projects": project, "statuses": "OPEN,REOPENED,CONFIRMED", "types": type})
        return response.body["total"]

    def get_issues_on_type_and_language(self, type, project, lang):
        response = unirest.get(self.url + "api/issues/search", headers=self.headers, auth=(self.user, self.psw), params={"projects": project, "statuses": "OPEN,REOPENED,CONFIRMED", "types": type, "languages": lang})
        return response.body["total"]

    def get_duplicated_blocks(self, project) :
        response = unirest.get(self.url + "api/measures/component", headers=self.headers, auth=(self.user, self.psw), params={"component": project, "metricKeys": "duplicated_blocks"})
        return response.body["component"]["measures"][0]["value"]

    def get_duplicated_density(self, project) :
        response = unirest.get(self.url + "api/measures/component", headers=self.headers, auth=(self.user, self.psw), params={"component": project, "metricKeys": "duplicated_lines_density"})
        return response.body["component"]["measures"][0]["value"]

    def get_violations(self, lang, type, project) : #podriamos devolver un codigo de error en caso de que hubiese problemas en vez de un array vacio
        try :
            r = list()
            sum = 0
            for language in lang :
                reglas = list()
                response = unirest.get(self.url + "api/issues/search", headers=self.headers, auth=(self.user, self.psw), params={"projectKeys": project, "facets": "rules", "types": type, "languages": language[3], "statuses": "OPEN,CONFIRMED,REOPENED"})
                rules = response.body["facets"][0]["values"]
                if len(rules) > 0 :
                    for i in range(len(rules)):
                        wget = unirest.get(self.url + "api/rules/show", headers=self.headers, auth=(self.user, self.psw), params={"key": rules[i]["val"]})
                        reglas.append((wget.body["rule"]["name"], resources.span_number_format(int(rules[i]["count"]))))
                        print round(float(i)/float(len(rules)), 3)*100, "%", language[0], "   \r",
                    sum += int(response.body["total"])
                    r.append((language[0], reglas))
            total = (type, resources.span_number_format(int(sum)), r)
            return total
        except :
            return list()

    def get_num_violations_by_lang(self, project, lang):
        sum = 0
        for type in ("VULNERABILITY,SECURITY_HOTSPOT", "BUG", "CODE_SMELL"):
            response = unirest.get(self.url + "api/issues/search", headers=self.headers, auth=(self.user, self.psw), params={"projectKeys": project, "facets": "rules", "types": type, "languages": lang, "statuses": "OPEN,CONFIRMED,REOPENED"})
            sum += len(response.body["facets"][0]["values"])
        return sum

    def get_num_rules_by_qprofile(self, qprofile):
        response = unirest.get(self.url + "api/rules/search", headers=self.headers, auth=(self.user, self.psw), params={"activation": "true", "qprofile": qprofile})
        return response.body["total"]

    def get_qprofile_key_by_lang(self, project, lang):
        response = unirest.get(self.url + "api/qualityprofiles/search", headers=self.headers, auth=(self.user, self.psw), params={"project": project, "language": lang})
        try :
            qprofile_key = response.body["profiles"][0]["key"]
            return qprofile_key
        except :
            return None

    def get_languages(self):
        response = unirest.get(self.url + "api/languages/list", headers=self.headers, auth=(self.user, self.psw))
        try:
            return response.body["languages"]
        except:
            return list()

    def get_rules_by_lang(self, lang) :
        response = unirest.get(self.url + "api/rules/search", headers=self.headers, auth=(self.user, self.psw), params={"languages": lang, "ps": "500"})
        try:
            return response.body["rules"]
        except:
            return list()

    def get_project_analyses_issues(self, project) :
        response = unirest.get(self.url + "api/measures/search_history", headers=self.headers, auth=(self.user, self.psw), params={"component": project, "metrics": "vulnerabilities,code_smells,bugs", "ps": "1000"})
        try:
            return response.body["measures"]
        except:
            return list()

    def get_project_analyses_duplicity(self, project):
        response = unirest.get(self.url + "api/measures/search_history", headers=self.headers, auth=(self.user, self.psw), params={"component": project, "metrics": "duplicated_blocks,duplicated_lines_density", "ps": "1000"})
        try:
            return response.body["measures"]
        except:
            return list()
