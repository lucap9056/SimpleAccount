<script lang="ts">
  import API from "../API";
  import { alertManager, Alert } from "../Alert/Struct";
  import Translations from "../Translations";
  import Loading from "../Loading/Main";
  import { router, Routes } from "../Router";
  import Validate from "../ValidateInput";

  let username = "";
  let email = "";
  let password = "";
  let retype_password = "";

  let usernameFailed = "";
  let emailFailed = "";
  let passwordFailed = "";
  let retype_passwordFailed = "";
  $: {
    usernameFailed = Validate.Username(username);
    emailFailed = Validate.Email(email);
    passwordFailed = Validate.Password(password);
    retype_passwordFailed = Validate.RetypePassword(password, retype_password);
  }

  async function handleSubmit() {
    const disable =
      usernameFailed !== "" || emailFailed !== "" || passwordFailed !== "";
    if (disable) {
      return;
    }

    const loading = Loading.Append();

    API.Register(username, email, password).then((res) => {
      if (res.success) {
        alertManager.Add(
          Translations.Get("register_email_verification"),
          Alert.Type.Alert,
          async () => router.Set(Routes.LOGIN),
          Translations.Get("register_confirm"),
        );
      } else {
        alertManager.Add(Translations.Get(res.error), Alert.Type.Error);
      }
      loading.Remove();
    });
  }
</script>

<div class="container">
  <form on:submit|preventDefault={handleSubmit}>
    <h2>{Translations.Get("register")}</h2>

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
      <label for="email">{Translations.Get("register_email")}</label>

      <div class="input" data-type="email">
        <input type="text" id="email" bind:value={email} required />
        {#if email != "" && emailFailed !== ""}
          <div class="input_alert">{emailFailed}</div>
        {/if}
      </div>
    </div>

    <div class="form-group">
      <label for="password">{Translations.Get("register_password")}</label>

      <div class="input" data-type="password">
        <input type="password" id="password" bind:value={password} required />
        {#if password != "" && passwordFailed !== ""}
          <div class="input_alert">{passwordFailed}</div>
        {/if}
        <div class="limit">
          {Translations.Get("register_password_limit")}
        </div>
      </div>
    </div>

    <div class="form-group">
      <label for="retype_password">
        {Translations.Get("register_retype_password")}
      </label>

      <div class="input" data-type="retype_password">
        <input
          type="password"
          id="retype_password"
          bind:value={retype_password}
          required
        />
        {#if retype_password != "" && retype_passwordFailed !== ""}
          <div class="input_alert">{retype_passwordFailed}</div>
        {/if}
      </div>
    </div>

    <div class="form-group">
      <button type="submit">{Translations.Get("register")}</button>
    </div>
  </form>
  <a href="/#login"><span>{Translations.Get("login")}</span></a>
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
