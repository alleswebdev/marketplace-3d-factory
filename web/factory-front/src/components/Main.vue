<template>
  <v-container>
    <v-tabs v-model="tab" align-tabs="left" color="deep-purple-accent-4">
      <v-card-title>Crazy Shop 3D</v-card-title>
      <v-tab value="wb">
        <v-badge color="error" :content="wbItems.length" floating>WB</v-badge>
      </v-tab>
      <v-tab value="ozon">
        <v-badge color="error" :content="ozonItems.length" floating>OZON</v-badge>
      </v-tab>
      <v-tab value="yandex">
        <v-badge color="error" :content="yandexItems.length" floating>Yandex</v-badge>
      </v-tab>
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
          <template v-slot:item.composite_items="{ item }">
            <v-row no-gutters style="height: 40px;">
              <v-col>
                <v-card-text>
                  {{ item.article }}
                </v-card-text>
              </v-col>
            </v-row>

            <v-row v-for="childrenItem in item.composite_items" no-gutters style="height: 40px;">
              <v-col>
                <v-card-text>
                  {{ childrenItem.name }}
                </v-card-text>
              </v-col>
              <v-col>
                <v-checkbox v-model="childrenItem.is_complete" @change="setChildrenCompleteV2(childrenItem)"
                            hide-details></v-checkbox>
              </v-col>
            </v-row>

          </template>
          <template v-slot:item.is_printing="{ item }">
            <v-checkbox v-model="item.is_printing" @change="setIsPrintingV2(item)"></v-checkbox>
          </template>
          <template v-slot:item.is_complete="{ item }">
            <v-btn @click="setCompleteV2(item)">{{ item.is_complete === true ? "Вернуть" : "Собрать" }}</v-btn>
          </template>
        </v-data-table>
      </v-window-item>

      <v-window-item value="ozon">
        <br>
        <v-tabs v-model="ozonSubTab" align-tabs="left" color="deep-purple-accent-4">
          <v-tab value="all">Все</v-tab>
          <v-tab v-for="(item, index) in groupedOzonItems" :value="index">{{index}}</v-tab>
        </v-tabs>
        <v-data-table
          :headers="ozonHeaders"
          :items="ozonSubTab === 'all' ? ozonItems : groupedOzonItems[ozonSubTab]"
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
          <template v-slot:item.composite_items="{ item }">
            <v-row no-gutters style="height: 40px;">
              <v-col>
                <v-card-text>
                  {{ item.article }}
                </v-card-text>
              </v-col>
            </v-row>

            <v-row v-for="childrenItem in item.composite_items" no-gutters style="height: 40px;">
              <v-col>
                <v-card-text>
                  {{ childrenItem.name }}
                </v-card-text>
              </v-col>
              <v-col>
                <v-checkbox v-model="childrenItem.is_complete" @change="setChildrenCompleteV2(childrenItem)"
                            hide-details></v-checkbox>
              </v-col>
            </v-row>

          </template>
          <template v-slot:item.is_printing="{ item }">
            <v-checkbox v-model="item.is_printing" @change="setIsPrintingV2(item)"></v-checkbox>
          </template>
          <template v-slot:item.is_complete="{ item }">
            <v-btn @click="setCompleteV2(item)">{{ item.is_complete === true ? "Вернуть" : "Собрать" }}</v-btn>
          </template>
        </v-data-table>
      </v-window-item>

      <v-window-item value="yandex">
        <br>
        <v-tabs v-model="yandexSubTab" align-tabs="left" color="deep-purple-accent-4">
          <v-tab value="all">Все</v-tab>
          <v-tab v-for="(item, index) in groupedYandexItems" :value="index">{{index}}</v-tab>
        </v-tabs>
        <v-data-table
          :headers="ozonHeaders"
          :items="yandexSubTab === 'all' ? yandexItems : groupedYandexItems[yandexSubTab]"
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
          <template v-slot:item.composite_items="{ item }">
            <v-row no-gutters style="height: 40px;">
              <v-col>
                <v-card-text>
                  {{ item.article }}
                </v-card-text>
              </v-col>
            </v-row>

            <v-row v-for="childrenItem in item.composite_items" no-gutters style="height: 40px;">
              <v-col>
                <v-card-text>
                  {{ childrenItem.name }}
                </v-card-text>
              </v-col>
              <v-col>
                <v-checkbox v-model="childrenItem.is_complete" @change="setChildrenCompleteV2(childrenItem)"
                            hide-details></v-checkbox>
              </v-col>
            </v-row>

          </template>
          <template v-slot:item.is_printing="{ item }">
            <v-checkbox v-model="item.is_printing" @change="setIsPrintingV2(item)"></v-checkbox>
          </template>
          <template v-slot:item.is_complete="{ item }">
            <v-btn @click="setCompleteV2(item)">{{ item.is_complete === true ? "Вернуть" : "Собрать" }}</v-btn>
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
      yandexItems: [],
      groupedOzonItems: [],
      groupedYandexItems: [],
      overlay: false,
      overlayScr: '',
      tab: null,
      ozonSubTab: null,
      yandexSubTab: null,
      appHost: "",
      headers: [
        {title: '', key: 'photo', sortable: false},
        {title: ' 🖨️', key: 'is_printing', sortable: false},
        // { title: '№ заказа', key: 'order_id', sortable: false},
        // { title: 'Артикул', key: 'article' , sortable: false},
        {title: 'Состав', key: 'composite_items', sortable: false},
        {title: 'Прошло времени', key: 'time_passed'},
        {title: 'Готов', key: 'is_complete', sortable: false}
      ],
      ozonHeaders: [
        {title: 'Готов', key: 'is_complete', sortable: false},
        {title: 'Номер отправления', key: 'info.order_number', sortable: false},
        {title: ' 🖨️', key: 'is_printing', sortable: false, align: 'center',},
        {title: 'Отгрузка', key: 'shipment_date'},
        {title: 'Состав', key: 'composite_items', sortable: false},
        {title: '', key: 'photo', sortable: false},
        {title: 'Количество', key: 'info.quantity', sortable: false},
      ],
    };
  },
  mounted() {
    this.appHost = process.env.MARKETPLACE_APP_HOST || "127.0.0.1"
    this.fetchWbItems();
    this.fetchOzonItems();
    this.fetchYandexItems();
  },
  created() {
    setInterval(this.fetchWbItems, 30000);
    setInterval(this.fetchOzonItems, 30000);
    setInterval(this.fetchYandexItems, 30000);
    const tabData = localStorage.getItem('tab');
    if (tabData) {
      this.tab = JSON.parse(tabData);
    }
  },
  watch: {
    tab(newValue, oldValue) {
      localStorage.setItem('tab', JSON.stringify(newValue));
    }
  },
  methods: {
    toggleOverlay(img) {
      this.overlay = true
      this.overlayScr = img
    },
    fetchItems() {
      this.fetchWbItems()
      this.fetchOzonItems()
      this.fetchYandexItems()
    },
    fetchWbItems() {
      if (this.isLoading) {
        return
      }
      this.isLoading = true

      axios.get(`/api/v2/list-queue?withParentComplete=${this.withCompleteParent}&marketplace=wb`)
        .then(response => {
          this.wbItems = response.data.items || [];
        })
        .catch(error => {
          console.error('Ошибка при получении данных:', error);
        });

      this.isLoading = false
    },
    groupByShipmentDate(response) {
      const groupedItems = {};

      response.items.forEach(item => {
        const shipmentDate = item.shipment_date;

        if (!groupedItems[shipmentDate]) {
          groupedItems[shipmentDate] = [];
        }

        groupedItems[shipmentDate].push(item);
      });

      return groupedItems;
    },
    fetchOzonItems() {
      if (this.isLoading) {
        return
      }
      this.isLoading = true

      axios.get(`/api/v2/list-queue?withParentComplete=${this.withCompleteParent}&marketplace=ozon`)
        .then(response => {
          this.ozonItems = response.data.items || [];
          this.groupedOzonItems = []
          if(this.ozonItems.length > 0){
            this.groupedOzonItems = this.groupByShipmentDate(response.data)
          }

        })
        .catch(error => {
          console.error('Ошибка при получении данных:', error);
        });

      this.isLoading = false
    },
    fetchYandexItems() {
      if (this.isLoading) {
        return
      }
      this.isLoading = true

      axios.get(`/api/v2/list-queue?withParentComplete=${this.withCompleteParent}&marketplace=yandex`)
        .then(response => {
          this.yandexItems = response.data.items || [];
          this.groupedYandexItems = []
          if(this.yandexItems.length > 0){
            this.groupedYandexItems = this.groupByShipmentDate(response.data)
          }

        })
        .catch(error => {
          console.error('Ошибка при получении данных:', error);
        });

      this.isLoading = false
    },
    // отдельные функции для кнопки собрать, потому что она не меняет стейт как чекбокс
    setCompleteV2(item) {
      axios.post('/api/v2/set-complete', {id: item.id, state: item.is_complete !== true})
        .then(response => {
          this.fetchItems()
        })
        .catch(error => {
          console.error('Ошибка при обновлении флага:', error);
        });
    },
    setIsCompleteV2(item) {
      axios.post('/api/v2/set-complete', {id: item.id, state: item.is_complete})
        .then(response => {
          this.fetchItems()
        })
        .catch(error => {
          console.error('Ошибка при обновлении флага:', error);
        });
    },
    setChildrenCompleteV2(item) {
      console.log(item)
      axios.post('/api/v2/set-children-complete', {id: item.id, state: item.is_complete})
        .then(response => {
          this.fetchItems()
        })
        .catch(error => {
          console.error('Ошибка при обновлении флага:', error);
        });
    },
    setIsPrintingV2(item) {
      axios.post('/api/v2/set-printing', {id: item.id, state: item.is_printing})
        .then(response => {
        })
        .catch(error => {
          console.error('Ошибка при обновлении флага:', error);
        });
    },
  },
};
</script>
