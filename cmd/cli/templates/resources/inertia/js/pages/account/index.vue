<template>
   <div class="mx-auto max-w-screen-xl px-4 py-16 sm:px-6 lg:px-8 lg:flex lg:h-screen lg:items-center">
        <div class="mx-auto max-w-lg">
            <Flash/>
            <Tag/>
            <LoginForm/>
        </div>
    </div>
</template>

<script setup>
    import Flash from "../../components/Flash.vue"
    import Tag from "../../components/Tag.vue"
    import LoginForm from "../../components/forms/LoginForm.vue"
</script>
