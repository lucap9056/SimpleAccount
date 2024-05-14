<script lang="ts">
	import Loading from "./Loading/Component.svelte";
	import Alert from "./Alert/Component.svelte";
	import RegisterForm from "./UI/RegisterForm.svelte";
	import LoginForm from "./UI/LoginForm.svelte";
	import Index from "./UI/Index.svelte";
	import { Routes, route, router } from "./Router";
	import Status from "./Status";
	Status.Login().then((login) => {
		switch (true) {
			case $route.INDEX:
				if (login) return;
				router.Set(Routes.LOGIN);
				break;
			case $route.LOGIN:
			case $route.REGISTER:
				if (!login) return;
				router.Set(Routes.INDEX);
				break;
		}
	});
</script>

<main>
	{#if $route.INDEX}
		<Index />
	{/if}
	{#if $route.LOGIN}
		<LoginForm />
	{/if}
	{#if $route.REGISTER}
		<RegisterForm />
	{/if}
	<Alert />
	<Loading />
</main>

<style>
	main {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		text-align: center;
		padding: 1em;
	}
</style>
