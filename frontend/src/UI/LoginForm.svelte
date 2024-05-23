<!-- RegisterForm.svelte -->
<script lang="ts">
    import { Alert, alertManager } from "../Alert/Struct";
    import Translations from "../Translations";
    import API from "../API";
    import loadings from "../Loading/Main";
    import Status from "../Status";
    import Textbox from "../Components/Textbox.svelte";
    import { Routes, router } from "../Router";

    let email = "";
    let password = "";

    async function handleSubmit() {
        const loading = loadings.Append();
        API.Login(email, password)
            .then(() => {
                Status.login.update(() => true);
                router.Set(Routes.INDEX);
                loading.Remove();
            })
            .catch((err) => {
                alertManager.Add($Translations[err], Alert.Type.Error);
                loading.Remove();
            });
    }
</script>

<div class="container">
    <form on:submit|preventDefault={handleSubmit}>
        <h2>{$Translations.login}</h2>
        <Textbox
            label={$Translations.login_email}
            name="email"
            bind:value={email}
        />

        <Textbox
            label={$Translations.login_password}
            name="password"
            password={true}
            bind:value={password}
        />

        <div class="form-group">
            <button type="submit">{$Translations.login}</button>
        </div>
    </form>
    <a href="/#register"><span>{$Translations.register}</span></a>
</div>

<style>
    .container {
        max-width: 400px;
        margin: 100px auto;
        background-color: #fff;
        padding: 20px;
        border-radius: 8px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    }

    h2 {
        text-align: center;
    }

    .form-group {
        margin: 20px;
    }

    button {
        width: 100%;
        padding: 10px;
        background-color: #007bff;
        color: #fff;
        border: none;
        border-radius: 5px;
        cursor: pointer;
    }

    button:hover {
        background-color: #0056b3;
    }
</style>
