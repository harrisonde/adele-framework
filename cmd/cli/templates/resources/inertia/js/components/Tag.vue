<template>

    <div v-if="type === 'banner'">
        <h1 class="bg-gradient-to-r from-green-300 via-blue-500 to-purple-600 bg-clip-text text-3xl font-extrabold text-transparent sm:text-5xl">
            {{ app.title }}
            <span class="sm:block"> {{ app.subTitle }}</span>
        </h1>

        <p class="mx-auto mt-4 max-w-xl sm:text-xl/relaxed">
            {{ app.paragraph }}
        </p>
    </div>
    <div v-else>
        <h1 class="text-center text-3xl font-bold text-indigo-600 sm:text-4xl">
            {{ app.title }}
        </h1>
        <h2 class="text-center text-2xl font-bold text-indigo-600 sm:text-3xl">
            {{  app.subTitle }}
        </h2>
        <p class="mx-auto mt-4 max-w-md text-center text-gray-500">
             {{ app.paragraph }}
        </p>
    </div>
</template>


<script setup>
    import { ref } from "vue"

    defineProps({
        type: String,
    })

    const app = ref({
        title: "Adele",
        subTitle: "Go, build awesomeness.",
        paragraph: "Lorem ipsum dolor sit amet consectetur, adipisicing elit. Nesciunt illotenetur fuga ducimus numquam ea!"
    })


</script>
