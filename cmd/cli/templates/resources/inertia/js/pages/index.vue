<template>
    <section class="bg-gray-900 text-white">
        <div class="mx-auto max-w-screen-xl px-4 py-32 lg:flex lg:h-screen lg:items-center">
            <div class="mx-auto max-w-3xl text-center">

                <Tag type="banner"/>

                <div class="mt-8 flex flex-wrap justify-center gap-4">
                    <Link
                        class="block w-full rounded border border-blue-600 bg-blue-600 px-12 py-3 text-sm font-medium text-white hover:bg-transparent hover:text-white focus:outline-none focus:ring active:text-opacity-75 sm:w-auto"
                        href="/account" method="get">Login</Link>
                    <Link
                        class="block w-full rounded border border-blue-600 px-12 py-3 text-sm font-medium text-white hover:bg-blue-600 focus:outline-none focus:ring active:bg-blue-500 sm:w-auto"
                        href="/account/register" method="get">Signup</Link>
                </div>
            </div>
        </div>
    </section>
</template>

<script setup>
    import { Link } from '@inertiajs/vue3'
    import Tag from "../components/Tag.vue"
</script>
