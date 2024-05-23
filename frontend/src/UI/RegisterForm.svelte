<script lang="ts">
  import API from "../API";
  import { alertManager, Alert } from "../Alert/Struct";
  import Translations from "../Translations";
  import Loading from "../Loading/Main";
  import { router, Routes } from "../Router";
  import Validate from "../ValidateInput";
  import Textbox from "../Components/Textbox.svelte";

  const data = {
    username: "",
    email: "",
    password: "",
    retype_password: "",
  };

  async function handleSubmit() {
    console.log(data);
    if (
      data.username == "" ||
      data.email == "" ||
      data.password == "" ||
      data.retype_password == ""
    ) {
      return;
    }

    const loading = Loading.Append();

    API.Register(data.username, data.email, data.password)
      .then(() => {
        alertManager.Add(
          $Translations.register_email_verification,
          Alert.Type.Alert,
          async () => router.Set(Routes.LOGIN),
          $Translations.register_confirm,
        );
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
    <h2>{$Translations.register}</h2>

    <Textbox
      label={$Translations.register_username}
      name="username"
      bind:value={data.username}
      validate={Validate.Username}
      hint={$Translations.register_username_limit}
    />

    <Textbox
      label={$Translations.register_email}
      name="email"
      bind:value={data.email}
      validate={Validate.Email}
    />

    <Textbox
      label={$Translations.register_password}
      name="password"
      password={true}
      bind:value={data.password}
      validate={Validate.Password}
      hint={$Translations.register_password_limit}
    />

    <Textbox
      label={$Translations.register_retype_password}
      name="retype_password"
      password={true}
      bind:value={data.retype_password}
      validate={(p) => Validate.RetypePassword(data.password, p)}
      hint={$Translations.register_password_limit}
    />

    <div class="form-group">
      <button type="submit">{$Translations.register}</button>
    </div>
  </form>
  <a href="/#login"><span>{$Translations.login}</span></a>
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
