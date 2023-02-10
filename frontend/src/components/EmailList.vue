<script>
export default {
  name: "EmailList",
  props: {
    emails: {
      type: Array,
      required: true,
    },
    total: {
      type: Number,
      required: true,
    },
  },
  data() {
    return {
      showModal: false,
      email: {},
      page: 1,
      perPage: 10,
      slicedEmails: [],
    };
  },
  mounted() {
    this.slicedEmails = this.sliceEmailsByPage();
  },
  methods: {
    onEmailClick(item) {
      this.email = item;
      this.toggleModal();
    },
    toggleModal() {
      this.showModal = !this.showModal;
    },
    sliceEmailsByPage() {
      const start = (this.page - 1) * this.perPage;
      const end = start + this.perPage;
      return this.emails.slice(start, end);
    },
    onNext() {
      this.page++;
      this.slicedEmails = this.sliceEmailsByPage();
    },
    onPrevious() {
      this.page--;
      this.slicedEmails = this.sliceEmailsByPage();
    },
  },
};
</script>

<template>
  <div class="border border-gray-300 rounded-lg px-6 py-3">
    <div
      class="grid grid-cols-12 gap-4 border-b border-b-gray-300 px-6 py-2 mb-1 w-full cursor-pointer"
      v-for="item in slicedEmails"
      :key="item?._id"
      @click="onEmailClick(item)"
    >
      <div class="truncate col-span-12 sm:col-span-4 font-medium">
        {{
          item?._source?.Subject && item?._source?.Subject !== " "
            ? item?._source?.Subject
            : "No Subject"
        }}
      </div>
      <div class="truncate col-span-12 sm:col-span-6">
        {{ item?._source?._body }}
      </div>
      <div class="col-span-12 sm:col-span-2 text-gray-400 text-right">
        {{ item?._source?.Date }}
      </div>
    </div>
    <div class="flex justify-center mt-6">
      <div class="flex space-x-2">
        <button
          v-if="page > 1"
          @click="onPrevious()"
          class="bg-white hover:bg-gray-100 text-gray-500 font-semibold py-1 px-3 border border-gray-400 rounded shadow"
        >
          Previous
        </button>
        <button
          v-if="page < Math.ceil(total.value / perPage)"
          @click="onNext()"
          class="bg-white hover:bg-gray-100 text-gray-500 font-semibold py-1 px-3 border border-gray-400 rounded shadow"
        >
          Next
        </button>
      </div>
    </div>
  </div>

  <!-- modal -->
  <div
    v-if="showModal"
    tabindex="1"
    style="background-color: rgba(255, 255, 255, 0.9)"
    class="fixed top-0 left-0 right-0 z-50 w-full p-4 overflow-x-hidden md:inset-0 h-modal"
  >
    <div
      class="ml-auto mr-auto relative w-full overflow-y-auto h-full max-w-2xl"
    >
      <!-- Modal content -->
      <div class="relative bg-white rounded-lg shadow border border-gray-300">
        <!-- Modal header -->
        <div class="flex items-start justify-between p-4 border-b rounded-t">
          <h3 class="truncate text-md font-semibold text-gray-900">
            {{
              item?._source?.Subject && item?._source?.Subject !== " "
                ? item?._source?.Subject
                : "No Subject"
            }}
          </h3>
          <button
            type="button"
            class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center dark:hover:bg-gray-600 dark:hover:text-white"
            @click="toggleModal()"
          >
            <svg
              aria-hidden="true"
              class="w-5 h-5"
              fill="currentColor"
              viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fill-rule="evenodd"
                d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                clip-rule="evenodd"
              ></path>
            </svg>
            <span class="sr-only">Close modal</span>
          </button>
        </div>
        <div class="flex justify-between p-4 border-b rounded-t">
          <div class="text-xs">
            <p class="mb-1">
              <span class="font-semibold">From:</span>
              {{ email?._source?.From }}
            </p>
            <p>
              <span class="font-semibold">To:</span> {{ email?._source?.To }}
            </p>
          </div>
          <div class="text-xs text-gray-400 text-right">
            {{ email?._source?.Date }}
          </div>
        </div>
        <!-- Modal body -->
        <div class="p-6 space-y-6">
          <p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
            {{ email?._source?._body }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
