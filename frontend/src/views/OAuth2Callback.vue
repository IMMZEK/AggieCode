<template>
    <div>
      <h1>Processing OAuth 2.0 Callback...</h1>
    </div>
  </template>
  
  <script>
  import { onMounted } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import oauth2Service from '../services/oauth2';
  import { useStore } from 'vuex';
  
  export default {
    setup() {
      const route = useRoute();
      const router = useRouter();
      const store = useStore();
  
      onMounted(async () => {
        const code = route.query.code;
        const state = route.query.state;
  
        if (!code || !state) {
          console.error('Invalid callback parameters');
          router.push('/');
          return;
        }
  
        try {
          const tokenData = await oauth2Service.exchangeCodeForToken(
            code,
            state
          );
  
          store.commit('oauth2/setOAuth2Tokens', tokenData);
  
          router.push('/gcs-authorized');
        } catch (error) {
          console.error('OAuth 2.0 flow failed:', error);
          router.push('/');
        }
      });
  
      return {};
    },
  };
  </script>