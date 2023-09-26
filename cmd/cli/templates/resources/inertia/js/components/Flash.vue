 <template>
    <div v-show="dismissed" v-if="flash && flash.type" role="alert" class="rounded-xl border border-gray-100 bg-white p-4 mb-4">
        <div class="flex items-start gap-4">
            <span v-if="flash.type == 'success'" class="text-green-600">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    class="h-6 w-6"
                >
                    <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                </svg>
            </span>

            <span v-else-if="flash.type == 'error'" class="text-red-600">

                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
                </svg>

            </span>

            <span v-else class="text-blue-600">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
                </svg>

            </span>

            <div class="flex-1">
            <strong class="block font-medium text-gray-900">{{ flash.title }}</strong>

            <p class="mt-1 text-sm text-gray-700">
                {{ flash.message }}
            </p>
            </div>

            <button v-on:click.prevent="dismiss" class="text-gray-500 transition hover:text-gray-600">
                <span class="sr-only">Dismiss</span>

                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    class="h-6 w-6"
                >
                    <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M6 18L18 6M6 6l12 12"
                    />
                </svg>
            </button>
        </div>
    </div>
</template>

<script setup>
    import { usePage } from '@inertiajs/vue3'
    import { computed, ref } from 'vue'

    const flash = computed( () => {
        let hasFlash = usePage().props.flash
        return hasFlash ? JSON.parse(hasFlash) : null
    })

    let dismissed = ref(1)

    const dismiss = () => {
        dismissed.value = dismissed.value ? 0 : 1
    }

</script>
