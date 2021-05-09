<template>
  <v-container>
    <v-row>
      <v-col v-if="templates.length !== 0" cols="3">
        <v-card tile>
            <v-list max-height="800" class="overflow-y-auto">
              <v-subheader>
                <v-text-field
                  v-model="searchQuery"
                  label="Search"
                  @input="filterTemplate"
                ></v-text-field>
              </v-subheader>
              <v-list-item-group color="primary">
                <v-list-item
                  v-for="(template, i) in temps"
                  :key="i"
                  v-on:click="updateSelectedTemplate(template)"
                >
                  <v-list-item-content>
                    <v-list-item-title v-text="template.name"></v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-card>
      </v-col>

      <v-col cols="9" v-if="selectedTemplate.content !== ''">
        <div class="ml-3">
          <div class="font-weight-bold body-1"> {{ selectedTemplate.name }} </div>
          <div v-if="selectedTemplate.compatible === false" class="red--text "> Resource API version not compatible with selected kubernetes version </div>
        </div>
        <code-viewer
            :readonly="true"
            :code="selectedTemplate.content"
        > </code-viewer>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import codeViewer from './CodeViewer'

export default {
    name: "ChartViewer",
    components: {
      codeViewer
    },
    props: [
      "templates"
    ],
    data() {
      return {
        selectedTemplate: {
          content: "",
          compatible: true
        },
        searchQuery: "",
        temps: this.templates
      }
    },
    methods: {
      updateSelectedTemplate(file) {
        this.selectedTemplate = file
      },
      filterTemplate() {
        console.log(this.searchQuery)
        const query = this.searchQuery
        this.temps = this.templates.filter(function (template) {
          return template.name.includes(query)
        })
      }
    }
  }
</script>