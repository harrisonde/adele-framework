<template>
    <div>
        <div v-if="errorBag" role="alert" class="mt-4 rounded border-s-4 border-red-500 bg-red-50 p-4">
            <strong class="block font-medium text-red-800"> Please correct the follow errors:</strong>

            <ul v-for="(value, key) in errorBag">
                <template v-for="(message, field) in value">
                    <li class="mt-2 ml-2 text-sm text-red-700">{{ message }}</li>
                </template>

            </ul>
        </div>


    <form @submit.prevent="submit" class="mb-0 mt-6 space-y-4 rounded-lg p-4 shadow-lg sm:p-6 lg:p-8">
        <p class="text-center text-lg font-medium">Sign up for your account</p>

        <div>
            <label for="first-name" class="sr-only">First Name</label>

            <div class="relative">
                <input
                    v-model="form.firstName"
                    type="text"
                    id="first-name"
                    class="w-full rounded-lg border-gray-200 p-4 pe-12 text-sm shadow-sm"
                    placeholder="Enter your first name"
                    autocomplete="given-name"
                />
            </div>

            <div v-if="errorBag && errorBag.Errors.firstName"><p class="mt-2 ml-2 text-sm text-red-700">{{ errorBag.Errors.firstName }}</p></div>
        </div>

        <div>
            <label for="last-name" class="sr-only">Last Name</label>

            <div class="relative">
                <input
                    v-model="form.lastName"
                    type="text"
                    id="last-name"
                    class="w-full rounded-lg border-gray-200 p-4 pe-12 text-sm shadow-sm"
                    placeholder="Enter your last name"
                    autocomplete="family-name"
                />
            </div>

            <div v-if="errorBag && errorBag.Errors.lastName"><p class="mt-2 ml-2 text-sm text-red-700">{{ errorBag.Errors.lastName }}</p></div>
        </div>

            <div>
            <label for="email" class="sr-only">Email</label>

            <div class="relative">
                <input
                    v-model="form.email"
                    type="email"
                    class="w-full rounded-lg border-gray-200 p-4 pe-12 text-sm shadow-sm"
                    placeholder="Enter your email address"
                    autocomplete="email"
                />
                <div v-if="errorBag && errorBag.Errors.email"><p class="mt-2 ml-2 text-sm text-red-700">{{ errorBag.Errors.email }}</p></div>
                <span class="absolute inset-y-0 end-0 grid place-content-center px-4">
                    <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-4 w-4 text-gray-400"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207"
                    />
                    </svg>
                </span>
            </div>
        </div>

        <div>
            <label for="password" class="sr-only">Password</label>

            <div class="relative">

                <input
                v-model="form.password"
                type="password"
                id="password"
                class="w-full rounded-lg border-gray-200 p-4 pe-12 text-sm shadow-sm"
                placeholder="Enter password"
                autocomplete="new-password"
            />

            <span class="absolute inset-y-0 end-0 grid place-content-center px-4">
                    <svg xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    class="h-4 w-4 text-gray-400">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88" />
                </svg>
            </span>
            </div>

            <div v-if="errorBag && errorBag.Errors.password"><p class="mt-2 ml-2 text-sm text-red-700">{{ errorBag.Errors.password }}</p></div>

        </div>

        <div>
            <label for="password-confirmation" class="sr-only">Confirm Password</label>

            <div class="relative">

                <input
                    v-model="form.passwordConfirmation"
                    type="password"
                    id="password-confirmation"
                    class="w-full rounded-lg border-gray-200 p-4 pe-12 text-sm shadow-sm"
                    placeholder="Confirm password"
                    autocomplete="new-password"
                />

                <span class="absolute inset-y-0 end-0 grid place-content-center px-4">
                    <svg xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    class="h-4 w-4 text-gray-400">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M3.98 8.223A10.477 10.477 0 001.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.45 10.45 0 0112 4.5c4.756 0 8.773 3.162 10.065 7.498a10.523 10.523 0 01-4.293 5.774M6.228 6.228L3 3m3.228 3.228l3.65 3.65m7.894 7.894L21 21m-3.228-3.228l-3.65-3.65m0 0a3 3 0 10-4.243-4.243m4.242 4.242L9.88 9.88" />
                    </svg>

                </span>

                <div v-if="errorBag && errorBag.Errors.passwordConfirmation"><p class="mt-2 ml-2 text-sm text-red-700">{{ errorBag.Errors.passwordConfirmation }}</p></div>
            </div>
        </div>


        <button
            type="submit"
            class="block w-full rounded-lg bg-indigo-600 px-5 py-3 text-sm font-medium text-white">
            Submit
        </button>

        <p class="text-center text-sm text-gray-500">
            Already have an account?
            <Link class="underline" href="/account" method="get">Login</Link>
        </p>
    </form>
        </div>
</template>

<script setup>
    import { Link, router} from '@inertiajs/vue3'
    import { reactive } from 'vue'

    const props = defineProps(['errorBag'])

    const form = reactive({
        email: null,
        password: null,
        firstName: null,
        lastName: null,
        password: null,
        passwordConfirmation: null
    })

    function submit(){
        router.post("/account/register", form, {errorBag: 'errorBag'})
    }

</script>
