// demo data
var treeData = {
    path: "",
    children: []
};

// define the tree-item component
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

// boot up the demo
var demo = new Vue({
    el: "#demo",
    data: {
        treeData: treeData,
        path: ""
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
        doAddPath() {
            const params = new URLSearchParams();
            params.append('path', this.path);

            axios.post('/api/addPath', params)
                .then(response => {
                    if (response.status != 200) {
                        throw new Error('something error');
                    } else {
                        this.doFetchSitemap();
                    }
                })
        }
    }
});