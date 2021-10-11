Vue.component("tree-item", {
    template: "#item-template",
    props: {
        item: Object
    },
    data: function () {
        return {
            isOpen: true
        };
    },
    computed: {
        isFolder: function () {
            return this.item.children && this.item.children.length;
        }
    },
    methods: {
        toggle: function () {
            if (this.isFolder) {
                this.isOpen = !this.isOpen;
            }
        },
    }
});

var tree = new Vue({
    el: "#tree",
    data: {
        treeData: {},
        url: ""
    },
    created: function () {
        this.doFetchSitemap();
    },
    methods: {
        doFetchSitemap() {
            axios.get('/api/sitemap')
                .then(response => {
                    if (response.status != 200) {
                        throw new Error('something error');
                    } else {
                        var resultSitemap = response.data;

                        this.treeData = resultSitemap;
                    }
                })
        },
        doCrawl() {
            const params = new URLSearchParams();
            params.append('url', this.url);

            axios.post('/api/crawl', params)
                .then(response => {
                    if (response.status != 200) {
                        throw new Error('something error');
                    } else {
                        this.doFetchSitemap();
                    }
                })
        },
        doFound() {
            axios.get('/api/found')
                .then(response => {
                    if (response.status != 200) {
                        throw new Error('something error');
                    } else {
                        var founditem = response.data;
                        console.log(founditem)

                    }
                })
        },
        doUpload() {
            let data = new FormData;
            data.append('sitemap', this.file)
            axios.post('/upload', data)
                .then(response => {
                    if (response.status != 200) {
                        throw new Error('something error');
                    } else {
                        this.doFetchSitemap();
                    }
                })
        },
        onUpload: function (event) {
            this.file = event.target.files[0]
        },
        doSort() {
            axios.get('/api/sort')
                .then(response => {
                    if (response.status != 200) {
                        throw new Error('something error');
                    } else {
                        this.doFetchSitemap();
                    }
                })
        },
        doReset() {
            axios.get('/api/reset')
                .then(response => {
                    if (response.status != 200) {
                        throw new Error('something error');
                    } else {
                        this.doFetchSitemap();
                    }
                })
        },
        doScan() {
            axios.get('/api/scan')
                .then(response => {
                    if (response.status != 200) {
                        throw new Error('something error');
                    } else {
                        var founditem = response.data;
                        console.log(founditem)
                    }
                })
        }
    }
});