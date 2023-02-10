<script>
import { getEmails } from "./services/zincsearch";
// components
import Spinner from "./components/Spinner.vue";
import EmailList from "./components/EmailList.vue";
import Navbar from "./components/Navbar.vue";
import Search from "./components/Search.vue";

export default {
  name: "App",
  data() {
    return {
      search: "",
      emails: [],
      total: 0,
      isLoading: false,
    };
  },
  async mounted() {
    await this.searchEmails();
  },
  components: {
    Spinner,
    EmailList,
    Navbar,
    Search,
  },
  methods: {
    async searchEmails() {
      this.isLoading = true;
      try {
        const data = await getEmails(this.search);
        const {
          hits: { hits, total },
        } = data;
        this.emails = hits;
        this.total = total;
        this.isLoading = false;
      } catch (error) {
        this.isLoading = false;
        throw error;
      }
    },
    async onSearchEmailsFromQuery(key) {
      if (key === "Enter") {
        await this.searchEmails();
      }
    },
  },
};
</script>

<template>
  <Navbar />
  <div class="container max-w-screen-lg mx-auto pt-10">
    <Search v-model="search" @searchEmails="onSearchEmailsFromQuery" />
  </div>
  <div class="container max-w-screen-lg mx-auto pt-10 w-full">
    <Spinner v-if="isLoading" />
    <EmailList v-else :emails="emails" :total="total" />
  </div>
</template>
