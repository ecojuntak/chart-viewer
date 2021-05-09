const axios = require('axios')

const baseURL = process.env.VUE_APP_API_SERVER_HOST

async function fetchRepos() {
    return await axios.get(baseURL + '/api/v1/repos')
}

async function fetchCharts(repoName) {
    return await axios.get(baseURL + '/api/v1/charts/' + repoName)
}

async function fetchChart(repoName, chartName, chartVersion, kubeVersion) {
    return await axios.get(baseURL + '/api/v1/charts/' + repoName + '/' + chartName + '/' + chartVersion + '?kube-version=' + kubeVersion)
}

async function renderManifest(repoName, chartName, chartVersion, values) {
    return await axios.post(baseURL + '/api/v1/charts/manifests/render/' + repoName + '/' + chartName + '/' + chartVersion, {'values': values})
    .then(res => res)
    .catch(err => err.response)
}

module.exports = {
    fetchRepos,
    fetchCharts,
    fetchChart,
    renderManifest
}
