import axios from "../utils/axios";

const formatQuerybody = (searchTerm) => ({
  query: {
    bool: {
      must: [
        {
          range: {
            "@timestamp": {
              gte: "2023-02-07T22:25:13.810Z",
              lt: "2024-02-07T22:55:13.810Z",
              format: "2006-01-02T15:04:05Z07:00",
            },
          },
        },
        searchTerm
          ? { query_string: { query: searchTerm } }
          : { match_all: {} },
      ],
    },
  },
  sort: ["-@timestamp"],
  from: 0,
  highlight: {
    pre_tags: ["<pre>"],
    post_tags: ["</pre>"],
    fields: {
      To: {
        pre_tags: [],
        post_tags: [],
      },
    },
  },
  size: 100,
  aggs: {
    histogram: {
      date_histogram: {
        field: "@timestamp",
        calendar_interval: "",
        fixed_interval: "30s",
      },
    },
  },
});

export const getEmails = async (searchTerm) => {
  try {
    const params = formatQuerybody(searchTerm);
    const { data } = await axios.post("/es/enron_mail/_search", params);
    return data;
  } catch (error) {
    console.log(error);
    throw error;
  }
};
