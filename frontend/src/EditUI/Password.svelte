<script lang="ts">
    import Translations from "../Translations";
    import Validate from "../ValidateInput";
    import API from "../API";
    import loadings from "../Loading/Main";
    import { Alert, alertManager } from "../Alert/Struct";
    import Textbox from "../Components/Textbox.svelte";
    export let Cancel: () => void;

    let currentPassword = "";
    let newPassword = "";
    let retypePassword = "";

    function handleSubmit() {
        if (
            currentPassword == "" ||
            newPassword == "" ||
            retypePassword == ""
        ) {
            return;
        }

        const loading = loadings.Append();

        API.ChangePassword(currentPassword, newPassword)
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

    function ValidateNewPassword(pswd: string): string {
        return Validate.Password(currentPassword, pswd);
    }

    function ValidateRetypePassword(pswd: string): string {
        return Validate.RetypePassword(newPassword, pswd);
    }
</script>

<form on:submit|preventDefault={handleSubmit}>
    <h2>{$Translations.index_change_password}</h2>

    <Textbox
        label={$Translations.index_change_current_password}
        name="password"
        bind:value={currentPassword}
        password={true}
        validate={Validate.Password}
    />

    <Textbox
        label={$Translations.index_change_new_password}
        name="new_password"
        bind:value={newPassword}
        password={true}
        validate={ValidateNewPassword}
        hint={$Translations.register_password_limit}
    />

    <Textbox
        label={$Translations.index_change_retype_password}
        name="new_password"
        bind:value={retypePassword}
        password={true}
        validate={ValidateRetypePassword}
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
        margin: 20px 0;
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
