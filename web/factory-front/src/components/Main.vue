<template>
  <v-container fluid>
    <v-tabs v-model="tab" align-tabs="left" color="deep-purple-accent-4">
      <v-card-title>Crazy Shop 3D</v-card-title>
      <v-tab value="wb">WB</v-tab>
      <v-tab value="ozon">OZON</v-tab>
    </v-tabs>
      <br>
        <v-row>
          <v-col class="py-2" cols="12">
            <v-btn-toggle
              v-model="withCompleteParent"
              color="deep-purple-accent-3"
            >
              <v-btn :value=false @click="fetchItems()">Очередь</v-btn>
              <v-btn :value=true @click="fetchItems()">Архив</v-btn>
            </v-btn-toggle>
          </v-col>
        </v-row>

        <v-window v-model="tab">
          <v-window-item value="wb">
            <v-data-table
              :headers="headers"
              :items="wbItems"
              :items-per-page="0"
              item-value="id"
              :hide-default-footer="true"
              height="calc(100vh - 180px)"
              fixed-header
            >
              <template #bottom></template>
              <template v-slot:top>
              </template>
              <template v-slot:item.photo="{ item }">
                <v-card class="my-2" elevation="2" width="100" rounded tile @click="toggleOverlay(item.photo)">
                  <v-img
                    :src="item.photo"
                    height="130"
                    width="100"
                    cover
                  ></v-img>
                </v-card>
              </template>
              <template v-slot:item.children="{ item }">
                <v-row no-gutters style="height: 40px;">
                  <v-col>
                    <v-card-text>
                      {{ item.article }}
                    </v-card-text>
                  </v-col>
                </v-row>

                <v-row v-for="item in item.children" no-gutters style="height: 40px;">
                  <v-col>
                      <v-card-text>
                        {{ item.article }}
                      </v-card-text>
                  </v-col>
                  <v-col>
                    <v-checkbox class="pa-0 ma-0" v-model="item.is_complete" @change="setIsComplete(item)" hide-details></v-checkbox>
                  </v-col>
                </v-row>

              </template>
              <template v-slot:item.is_printing="{ item }">
                <v-checkbox v-model="item.is_printing" @change="setIsPrinting(item)"></v-checkbox>
              </template>
              <template v-slot:item.is_complete="{ item }">
                <v-btn @click="setIsComplete(item)">Собрать</v-btn>
              </template>
            </v-data-table>
          </v-window-item>

          <v-window-item value="ozon">
            <v-data-table
              :headers="headers"
              :items="ozonItems"
              :items-per-page="0"
              item-value="id"
              :hide-default-footer="true"
              height="calc(100vh - 180px)"
              fixed-header
            >
              <template #bottom></template>
              <template v-slot:top>
              </template>
              <template v-slot:item.photo="{ item }">
                <v-card class="my-2" elevation="2" width="100" rounded tile @click="toggleOverlay(item.photo)">
                  <v-img
                    :src="item.photo"
                    height="130"
                    width="100"
                    cover
                  ></v-img>
                </v-card>
              </template>
              <template v-slot:item.children="{ item }">
                <v-row no-gutters style="height: 40px;">
                  <v-col>
                    <v-card-text>
                      {{ item.article }}
                    </v-card-text>
                  </v-col>
                </v-row>

                <v-row v-for="item in item.children" no-gutters style="height: 40px;">
                  <v-col>
                    <v-card-text>
                      {{ item.article }}
                    </v-card-text>
                  </v-col>
                  <v-col>
                    <v-checkbox class="pa-0 ma-0" v-model="item.is_complete" @change="setIsComplete(item)" hide-details></v-checkbox>
                  </v-col>
                </v-row>

              </template>
              <template v-slot:item.is_printing="{ item }">
                <v-checkbox v-model="item.is_printing" @change="setIsPrinting(item)"></v-checkbox>
              </template>
              <template v-slot:item.is_complete="{ item }">
                <v-btn @click="setIsComplete(item)">Собрать</v-btn>
              </template>
            </v-data-table>
          </v-window-item>
        </v-window>
  </v-container>
  <v-dialog v-model="overlay" max-width="500">
    <v-card>
      <v-img :src="overlayScr"></v-img>
    </v-card>
  </v-dialog>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      withCompleteParent: false,
      isLoading: false,
      expanded: [],
      wbItems: [],
      ozonItems: [],
      overlay: false,
      overlayScr: '',
      tab: null,
      appHost: "",
      headers: [
        {title: '', key: 'photo', sortable: false},
        {title: ' 🖨️', key: 'is_printing', sortable: false},
        // { title: '№ заказа', key: 'order_id', sortable: false},
        // { title: 'Артикул', key: 'article' , sortable: false},
        {title: 'Состав', key: 'children'},
        {title: 'Прошло времени', key: 'time_passed'},
        {title: 'Готов', key: 'is_complete', sortable: false}
      ],
    };
  },
  mounted() {
    this.appHost = process.env.MARKETPLACE_APP_HOST || "127.0.0.1"
    this.fetchWbItems();
    this.fetchOzonItems();
  },
  created() {
    setInterval(this.fetchWbItems, 30000);
    setInterval(this.fetchOzonItems, 30000);
  },
  methods: {
    toggleOverlay(img) {
      this.overlay = true
      this.overlayScr = img
    },
    fetchItems() {
      this.fetchWbItems()
      this.fetchOzonItems()
    },
    fetchWbItems() {
      if (this.isLoading) {
        return
      }
      this.isLoading = true

      axios.get(`/api/list-queue?withParentComplete=${this.withCompleteParent}&withChildrenComplete=true&marketplace=wb`)
        .then(response => {
          this.wbItems = response.data.items;
        })
        .catch(error => {
          console.error('Ошибка при получении данных:', error);
        });

      this.isLoading = false
    },
    fetchOzonItems() {
      if (this.isLoading) {
        return
      }
      this.isLoading = true

      axios.get(`/api/list-queue?withParentComplete=${this.withCompleteParent}&withChildrenComplete=true&marketplace=ozon`)
        .then(response => {
          this.ozonItems = response.data.items;
        })
        .catch(error => {
          console.error('Ошибка при получении данных:', error);
        });

      this.isLoading = false
    },
    setIsComplete(item) {
      item.is_complete = !item.is_complete
      axios.post('/api/set-complete', {id: item.id, state: item.is_complete})
        .then(response => {
          this.fetchWbItems()
        })
        .catch(error => {
          console.error('Ошибка при обновлении флага:', error);
        });
    },
    setIsPrinting(item) {
      axios.post('/api/set-printing', {id: item.id, state: item.is_printing})
        .then(response => {
        })
        .catch(error => {
          console.error('Ошибка при обновлении флага:', error);
        });
    },
  },
};
</script>
