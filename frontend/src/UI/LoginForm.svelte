<!-- RegisterForm.svelte -->
<script lang="ts">
    import { Alert, alertManager } from "../Alert/Struct";
    import Translations from "../Translations";
    import { Post, Path } from "../API";
    import loadings from "../Loading/Main";
    import { Routes, router } from "../Router";

    let username = "";
    let password = "";

    async function handleSubmit() {
        const loading = loadings.Append();
        Post(Path.LOGIN, {
            username,
            password,
        }).then((res) => {
            if (res.success) {
                router.Set(Routes.INDEX);
            } else {
                alertManager.Add(Translations.Get(res.error), Alert.Type.Error);
            }
            loading.Remove();
        });
    }
</script>

<div class="container">
    <form on:submit|preventDefault={handleSubmit}>
        <h2>{Translations.Get("login")}</h2>
        <div class="form-group">
            <label for="username">{Translations.Get("login_username")}</label>
            <input type="text" id="username" bind:value={username} required />
        </div>
        <div class="form-group">
            <label for="password">{Translations.Get("login_password")}</label>
            <input
                type="password"
                id="password"
                bind:value={password}
                required
            />
        </div>
        <div class="form-group">
            <button type="submit">{Translations.Get("login")}</button>
        </div>
    </form>
    <a href="/#register"><span>{Translations.Get("register")}</span></a>
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
        margin-bottom: 20px;
    }

    label {
        display: block;
        margin-bottom: 5px;
    }

    input[type="text"],
    input[type="password"] {
        width: 100%;
        padding: 10px;
        border-radius: 5px;
        border: 1px solid #ccc;
        box-sizing: border-box;
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
