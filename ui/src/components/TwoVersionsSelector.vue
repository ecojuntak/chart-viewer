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
          @change="fetchCharts"
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

    <v-row v-if="selectedChart != '' && (firstSelectedVersion == '' || secondSelectedVersion == '')" class="d-flex flex-row-reverse">
      <v-alert
        outlined
        dense
        type="warning"
        text
      >
        Please choose two versions 
      </v-alert>
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

  export default {
    name: 'TwoVersionsSelector',
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
      this.fetchRepos()
    },
    methods: {
      async fetchRepos() {
        const response = await api.fetchRepos()
        this.repos = response.data
      },
      async fetchCharts() {
        this.resetState()
        this.selectedChart = ""

        const response = await api.fetchCharts(this.selectedRepo)
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
        const firstChart = await api.fetchChart(this.selectedRepo, this.selectedChart, this.firstSelectedVersion)
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
        const secondChart = await api.fetchChart(this.selectedRepo, this.selectedChart, this.secondSelectedVersion)
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
        this.firstSelectedVersion = ""
        this.secondSelectedVersion = ""
      }
    },
    watch: {
      selectedRepo() {
        this.$emit("selectedRepoChanged", this.selectedRepo)
      },
      selectedChart() {
        this.$emit("selectedChartChanged", this.selectedChart)
      },
      firstTemplates() {
        this.$emit("firstTemplatesChanged", this.firstTemplates);
      },
      secondTemplates() {
        this.$emit("secondTemplatesChanged", this.secondTemplates);
      },
      firstValues() {
        this.$emit("firstValuesChanged", this.firstValues)
      },
      secondValues() {
        this.$emit("secondValuesChanged", this.secondValues)
      },
      firstSelectedVersion() {
        this.$emit("firstSelectedVersionChanged", this.firstSelectedVersion)
      },
      secondSelectedVersion() {
        this.$emit("secondSelectedVersionChanged", this.secondSelectedVersion)
      }
    }
  }
</script>