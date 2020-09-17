const axios = require('axios')

const baseURL = process.env.VUE_APP_API_SERVER_HOST

async function fetchRepoList() {
    return await axios.get(baseURL + '/repos')
}

async function fetchChartList(repoName) {
    return await axios.get(baseURL + '/repos/' + repoName)
}

async function fetchChartDetail(repoName, chartName, chartVersion) {
    return await axios.get(baseURL + '/repos/' + repoName + '/' + chartName + '/' + chartVersion)
}

async function fetchManifest(repoName, chartName, chartVersion, values) {
    return await axios.get(baseURL + '/repos/' + repoName + '/' + chartName + '/' + chartVersion + '/render?values=' + values)
}

module.exports = {
    fetchRepoList,
    fetchChartList,
    fetchChartDetail,
    fetchManifest
}
