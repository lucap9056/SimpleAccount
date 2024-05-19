<script lang="ts">
    import { visible } from "./Main";
</script>

{#if $visible}
    <div id="loading">
        <img id="loading_image" src="/assets/loading.svg" alt="" />
    </div>
{/if}

<style>
    #loading {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        display: flex;
        justify-content: center;
        align-items: center;
    }

    #loading_image {
        width: 256px;
        height: 256px;
        filter: drop-shadow(0 10px 4px #000a);
    }
</style>
