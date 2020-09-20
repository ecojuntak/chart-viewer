<template>
  <v-app id="app">
    <v-navigation-drawer v-model="drawer" app clipped>
      <v-list>
        <router-link to="/" class="no-style">
          <v-list-item link>
            <v-list-item-action>
              <v-icon>mdi-text-box-search-outline</v-icon>
            </v-list-item-action>
            <v-list-item-content>
              <v-list-item-title>
                Inspect 
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </router-link>

        <router-link to="/compare" class="no-style">
          <v-list-item link>
            <v-list-item-action>
              <v-icon>mdi-file-compare</v-icon>
            </v-list-item-action>
            <v-list-item-content>
              Compare 
            </v-list-item-content>
          </v-list-item>
        </router-link>

        <router-link to="/render-manifest" class="no-style">
          <v-list-item link>
            <v-list-item-action>
              <v-icon>mdi-tray-full</v-icon>
            </v-list-item-action>
            <v-list-item-content>
              Render Manifest 
            </v-list-item-content>
          </v-list-item>
        </router-link>

        <router-link to="/compare-manifest" class="no-style">
          <v-list-item link>
            <v-list-item-action>
              <v-icon>mdi-select-compare</v-icon>
            </v-list-item-action>
            <v-list-item-content>
              Compare Manifest 
            </v-list-item-content>
          </v-list-item>
        </router-link>

      </v-list>
    </v-navigation-drawer>

    <v-app-bar
      app
      color="primary"
      dark
      clipped-left
    >
      <div class="d-flex align-center">
        <h1>Chart Viewer</h1>
      </div>

      <v-spacer></v-spacer>

      <v-btn
        href="https://github.com/ecojuntak/chart-viewer"
        target="_blank"
        text
      >
        <v-icon>mdi-github</v-icon>
      </v-btn>
    </v-app-bar>

    <v-main>
      <router-view> </router-view>
    </v-main>
  </v-app>
</template>

<script>

export default {
  name: 'App',
  data: () => ({
    drawer: null,
  }),
};
</script>

<style scoped>
  .no-style {
    text-decoration: none;
  }
</style>
