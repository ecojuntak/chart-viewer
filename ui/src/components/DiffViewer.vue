<template>
  <v-container>
    <v-row>
      <v-col cols="3">
        <v-card tile>
            <v-list max-height="800" class="overflow-y-auto">
              <v-subheader>
                <v-text-field
                  v-model="searchQuery"
                  label="Search"
                  @input="filterTemplate"
                ></v-text-field>
              </v-subheader>
              <v-list-item-group color="primary" v-if="filteredTemps.length != 0">
                <v-list-item
                  v-for="(template, i) in filteredTemps"
                  :key="i"
                  @click="updateSelectedTemplate(template)"
                >
                  <v-list-item-content v-if="template.changed" class="yellow lighten-3">
                    <v-list-item-title v-text="template.name"></v-list-item-title>
                  </v-list-item-content>
                  <v-list-item-content v-else>
                    <v-list-item-title v-text="template.name"></v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-card>
      </v-col>

      <v-col cols="9">
        <div v-show="selectedTemplate.changed">
          <v-btn small color="primary" class="ma-2" @click="switchView()" v-if="selectedTemplate.changed"> Switch view </v-btn>
          <code-diff
            :old-string="selectedTemplate.firstContent" 
            :new-string="selectedTemplate.secondContent" 
            :context="20" 
            :outputFormat="format"
            :fileName="selectedTemplate.name"
          />
        </div>
        <div v-if="selectedTemplate.name != null && !selectedTemplate.changed">
          <code-viewer :code="selectedTemplate.firstContent" :message="selectedTemplate.name"> </code-viewer>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import codeDiff from 'vue-code-diff'
import codeViewer from './CodeViewer'

export default {
  components: { codeDiff, codeViewer },
  name: 'DiffViewer',
  props: [
    'firstTemplates',
    'secondTemplates',
  ],
  data () {
    return {
      searchQuery: "",
      temps: this.getMostCompletedTemplate(),
      filteredTemps: [],
      format: "side-by-side",
      selectedTemplate: {},
      delayed: false,
      emptyRenderMessage: '# This file empty because the file not expected to be rendered'
    }
  },
  methods: {
    filterTemplate() {
      const query = this.searchQuery
      const temps = this.temps.filter(function(temp){
                return temp.name.includes(query)
              });
            
      this.filteredTemps = temps
    },
    updateSelectedTemplate(file) {
      this.selectedTemplate = file
    },
    getMostCompletedTemplate() {
      if(this.firstTemplates.length >= this.secondTemplates.length) {
        return this.firstTemplates
      }
      return this.secondTemplates
    },
    switchView() {
      if(this.format === 'side-by-side') {
        this.format = 'line-by-line'
      } else {
        this.format = 'side-by-side'
      }
    },
    mergeTemplates() {
      var mergedTemplates = []
      const mostCompletedTemplate = this.getMostCompletedTemplate()

      for(var i=0; i < mostCompletedTemplate.length; i++) {
        const anchorTemplate = mostCompletedTemplate[i]
        const templateOne = this.firstTemplates.find(temp => temp.name === anchorTemplate.name)
        const templateTwo = this.secondTemplates.find(temp => temp.name === anchorTemplate.name)

        var contentOne = this.emptyRenderMessage
        var contentTwo = this.emptyRenderMessage
        if(templateOne != undefined) {
          contentOne = templateOne.content
        }

        if(templateTwo != undefined) {
          contentTwo = templateTwo.content
        }

        const mergedTemplate = {
          name: anchorTemplate.name,
          firstContent: contentOne,
          secondContent: contentTwo,
          changed: contentOne !== contentTwo
        }

        mergedTemplates.push(mergedTemplate)
      }

      return mergedTemplates.sort((t1, t2) => {
        return t1.changed !== t2.changed
      })
    }
  },
  async mounted() {
    this.filteredTemps = this.mergeTemplates()
  }
}
</script>

<style scoped>

</style>