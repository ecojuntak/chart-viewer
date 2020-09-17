<template>
  <v-container>
    <v-row align="center">
      <v-col class="d-flex" cols="3">
        <v-autocomplete
          :items="repos"
          label="Repo Name"
          :item-text="contructChartName"
          item-value="name"
          v-model="selectedRepo"
          @change="fetchChartList"
        ></v-autocomplete>
      </v-col>

      <v-col class="d-flex" cols="3">
        <v-autocomplete
          :items="charts"
          label="Chart Name"
          item-text="name"
          v-model="selectedChart"
          @change="fetchVersionList"
        ></v-autocomplete>
      </v-col>

      <v-col class="d-flex" cols="3">
        <v-autocomplete
          :items="versions"
          label="Chart Version"
          v-model="firstSelectedVersion"
          @change="fetchFirstChartDetail"
        ></v-autocomplete>
      </v-col>

      <v-col class="d-flex" cols="3">
        <v-autocomplete
          :items="versions"
          label="Chart Version"
          v-model="secondSelectedVersion"
          @change="fetchSecondChartDetail"
        ></v-autocomplete>
      </v-col>
    </v-row>

    <v-row v-if="selectedChart != '' && (firstSelectedVersion == '' || secondSelectedVersion == '')">
      <v-alert
        outlined
        type="warning"
        text
      >
        Please choose two versions 
      </v-alert>
    </v-row>

    <v-row v-if="firstTemplates.length != 0 && secondTemplates.length != 0">
      <diff-viewer :firstTemplates="firstTemplates" :secondTemplates="secondTemplates"> </diff-viewer>
    </v-row>

    <v-row v-if="progressing">
      <v-progress-linear
        indeterminate
        color="blue darken-2"
      ></v-progress-linear>
    </v-row>
  </v-container>
</template>

<script>
  import api from '../api/api'
  import yaml from 'json-to-pretty-yaml'
  import diffViewer from '../components/DiffViewer'

  export default {
    name: 'Compare',
    components: {
      diffViewer,
    },
    data () {
      return {
        repos: [],
        selectedRepo: "",
        charts: [],
        selectedChart: "",
        versions: [],
        firstSelectedVersion: "",
        secondSelectedVersion: "",
        firstValues: "",
        secondValues: "",
        firstTemplates: [],
        secondTemplates: [],
        progressing: false
      }
    },
    mounted() {
      this.fetchRepoList()
    },
    methods: {
      async fetchRepoList() {
        const response = await api.fetchRepoList()
        this.repos = response.data
      },
      async fetchChartList() {
        this.resetState()

        const response = await api.fetchChartList(this.selectedRepo)
        this.charts = response.data
      },
      fetchVersionList() {
        this.resetState()

        for(let i=0; i < this.charts.length; i++) {
          if(this.charts[i].name === this.selectedChart) {
            this.versions = this.charts[i].versions
            break
          }
        }
      },
      async fetchFirstChartDetail() {
        this.firstValues = ""
        this.firstTemplates = []

        this.progressing = true
        const firstChart = await api.fetchChartDetail(this.selectedRepo, this.selectedChart, this.firstSelectedVersion)
        this.progressing = false
        
        this.firstValues = yaml.stringify(firstChart.data.values)
        const firstTemps = firstChart.data.templates
        this.firstTemplates = this.simplifyTemplateName(firstTemps)
        this.firstTemplates.push({
          name: 'values.yaml',
          content: this.firstValues
        })
      },
       async fetchSecondChartDetail() {
        this.secondValues = ""
        this.secondTemplates = []

        this.progressing = true
        const secondChart = await api.fetchChartDetail(this.selectedRepo, this.selectedChart, this.secondSelectedVersion)
        this.progressing = false
        
        this.secondValues = yaml.stringify(secondChart.data.values)
        const secondTemps = secondChart.data.templates
        this.secondTemplates = this.simplifyTemplateName(secondTemps)
        this.secondTemplates.push({
          name: 'values.yaml',
          content: this.secondValues  
        })
      },
      simplifyTemplateName(templates) {
        var temps = []
        templates.forEach((template) => {
          const newName = template.name.replace("templates/", "")
          temps.push({
            name: newName,
            content: template.content
          })
        })

        return temps
      },
      contructChartName(repo) {
        return repo.name + " (" + repo.url + ")"
      },
      resetState() {
        this.versions = []
        this.firstTemplates = []
        this.secondTemplates = []
      }
    }  
  }
</script>