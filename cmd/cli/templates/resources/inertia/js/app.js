import { createApp, h } from 'vue'
import { createInertiaApp } from '@inertiajs/vue3'
import axios from 'axios'
import '../css/tailwind.css'

axios.interceptors.request.use((config) => {
    config.headers['X-CSRF-Token'] = document.querySelector('meta[name="csrf-token"]').content
    return config;
});

createInertiaApp({
    resolve: name => {
        try {
            const pages = import.meta.glob('./pages/**/*.vue', { eager: true })
            let hasTemplate = () => {
                for (const [key, value] of Object.entries(pages)) {
                    let templatePath = key.replace("./pages/", "").replace(/\.[^/.]+$/, "")
                    if (templatePath == name) {
                        return true
                    }
                }
            }
            if (hasTemplate() == undefined) {
                return pages[`./pages/404.vue`]
            }
            return pages[`./pages/${name}.vue`]
        } catch (e) {
            console.log(e)
        }
    },
    setup({ el, App, props, plugin }) {
        createApp({ render: () => h(App, props) })
            .use(plugin)
            .mount(el)
    }
})
