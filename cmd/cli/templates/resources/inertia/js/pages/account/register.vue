<template>
   <div class="mx-auto max-w-screen-xl px-4 py-16 sm:px-6 lg:px-8 lg:flex lg:h-screen lg:items-center">

        <div class="mx-auto max-w-lg">
            <Tag/>
            <RegistrationForm :errorBag="errorBag"/>
        </div>
    </div>
</template>

<script setup>
    import Tag from "../../components/Tag.vue"
    import RegistrationForm from '../../components/forms/RegistrationForm.vue';

    const props = defineProps(['errorBag'])
</script>
