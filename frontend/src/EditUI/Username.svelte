<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import Translations from "../Translations";
    import Validate from "../ValidateInput";
    import API from "../API";
    import loadings from "../Loading/Main";
    import { Alert, alertManager } from "../Alert/Struct";
    const dispatch = createEventDispatcher();
    function Cancel() {
        dispatch("cancel");
    }

    let username = "";
    let password = "";

    let usernameFailed = "";
    $: {
        usernameFailed = Validate.Username(username);
    }

    function handleSubmit() {
        if (usernameFailed != "") {
            return;
        }

        const loading = loadings.Append();

        API.ChangeUsername(username, password).then((res) => {
            if (res.success) {
                Cancel();
                alertManager.Add(
                    Translations.Get("index_change_success"),
                    Alert.Type.Normal,
                );
            } else {
                alertManager.Add(Translations.Get(res.error), Alert.Type.Error);
            }
            loading.Remove();
        });
    }
</script>

<form on:submit|preventDefault={handleSubmit}>
    <h2>{Translations.Get("index_change_username")}</h2>

    <div class="form-group">
        <label for="username">{Translations.Get("register_username")}</label>
        <div class="input" data-type="username">
            <input type="text" id="username" bind:value={username} required />
            {#if username != "" && usernameFailed !== ""}
                <div class="input_alert">{usernameFailed}</div>
            {/if}
            <div class="limit">
                {Translations.Get("register_username_limit")}
            </div>
        </div>
    </div>

    <div class="form-group">
        <label for="password"
            >{Translations.Get("index_current_password")}</label
        >

        <div class="input" data-type="password">
            <input
                type="password"
                id="password"
                bind:value={password}
                required
            />
        </div>
    </div>
    <div class="form-group options">
        <button type="button" on:click={Cancel}>
            {Translations.Get("index_edit_cancel")}
        </button>
        <button type="submit">
            {Translations.Get("index_edit_confirm")}
        </button>
    </div>
</form>

<style>
    h2 {
        text-align: center;
    }

    .form-group {
        margin-bottom: 20px;
    }

    .options {
        display: flex;
        justify-content: flex-end;
    }

    label {
        display: block;
        margin-bottom: 5px;
    }

    .input {
        position: relative;
    }

    input {
        position: relative;
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
        width: 80px;
        padding: 10px;
        background-color: #007bff;
        color: #fff;
        border: none;
        border-radius: 5px;
        cursor: pointer;
        margin-left: 20px;
    }

    button:hover {
        background-color: #0056b3;
    }

    .limit {
        font-size: 12px;
        color: #999;
        user-select: none;
        text-align: left;
    }

    .input_alert {
        position: absolute;
        top: 100%;
        white-space: nowrap;
        height: 32px;
        line-height: 32px;
        padding: 0 10px;
        background-color: #c46;
        border-radius: 5px;
        color: white;
        margin: 4px;
    }
</style>
