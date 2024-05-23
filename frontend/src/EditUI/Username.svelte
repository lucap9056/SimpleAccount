<script lang="ts">
    import Translations from "../Translations";
    import Validate from "../ValidateInput";
    import API from "../API";
    import loadings from "../Loading/Main";
    import { Alert, alertManager } from "../Alert/Struct";
    import Textbox from "../Components/Textbox.svelte";
    export let Cancel: () => void;

    let username = "";
    let password = "";
    function handleSubmit() {
        if (username == "" || password == "") {
            return;
        }

        const loading = loadings.Append();

        API.ChangeUsername(username, password)
            .then(() => {
                Cancel();
                alertManager.Add(
                    $Translations.index_change_success,
                    Alert.Type.Normal,
                );
                loading.Remove();
            })
            .catch((err) => {
                alertManager.Add($Translations[err], Alert.Type.Error);
                loading.Remove();
            });
    }
</script>

<form on:submit|preventDefault={handleSubmit}>
    <h2>{$Translations.index_change_username}</h2>

    <Textbox
        label={$Translations.register_username}
        name="username"
        bind:value={username}
        validate={Validate.Username}
        hint={$Translations.register_username_limit}
    />

    <Textbox
        label={$Translations.index_current_password}
        name="password"
        bind:value={password}
        hint={$Translations.register_username_limit}
        validate={Validate.Password}
        password={true}
    />

    <div class="form-group options">
        <button type="button" on:click={Cancel}>
            {$Translations.index_edit_cancel}
        </button>
        <button type="submit">
            {$Translations.index_edit_confirm}
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
</style>
