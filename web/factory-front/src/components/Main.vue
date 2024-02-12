<template>
  <v-data-table
    :headers="headers"
    :items="items"
    :items-per-page="0"
    item-value="id"
    :hide-default-footer="true"
    show-expand
    height="calc(100vh - 180px)"
    fixed-header
    v-model:expanded="expanded"
  >
    <template #bottom></template>
    <template v-slot:top>
      <v-toolbar flat>
        <v-toolbar-title>Очередь на печать</v-toolbar-title>
        <v-checkbox v-model="withCompleteParent" @change="fetchItems">Показывать завершённые</v-checkbox>
      </v-toolbar>
    </template>
    <template v-slot:expanded-row="{ columns, item }">
        <tr v-for="item in item.children">
          <td></td>
          <td>
            <v-card class="my-2" elevation="2"   width="75" rounded tile>
              <v-img :src="item.photo" height="100" width="75"></v-img>
            </v-card></td>
          <td>{{ item.article }}</td>
          <td>{{ item.time_passed }}</td>
          <td>{{ item.color }}</td>
          <td>{{ item.size }}</td>
          <td>{{ item.marketplace }}</td>
          <td>  <v-checkbox v-model="item.is_printing" @change="setIsPrinting(item)"></v-checkbox> </td>
          <td>  <v-checkbox v-model="item.is_complete" @change="setIsComplete(item)"></v-checkbox> </td>
        </tr>
    </template>
    <template v-slot:item.photo="{ item }">
      <v-card class="my-2" elevation="2"   width="75" rounded tile>
        <v-img
          :src="item.photo"
          height="100"
          width="75"
        ></v-img>
      </v-card>
    </template>
    <template v-slot:item.is_printing="{ item }"><v-checkbox v-model="item.is_printing" @change="setIsPrinting(item)"></v-checkbox></template>
    <template v-slot:item.is_complete="{ item }"><v-checkbox v-model="item.is_complete" @change="setIsComplete(item)"></v-checkbox></template>
  </v-data-table>
</template>

<script>
import axios from 'axios';
import {th} from "vuetify/locale";

export default {
  data() {
    return {
      withCompleteParent:false,
      isLoading:false,
      expanded: [],
      items: [],
      headers: [
        { title: '', key: 'photo', sortable: false},
        { title: '№ заказа', key: 'order_id', sortable: false},
        { title: 'Артикул', key: 'article' , sortable: false},
        { title: 'Прошло времени', key: 'time_passed' },
        { title: 'Печатается', key: 'is_printing' , sortable: false},
        { title: 'Готов', key: 'is_complete' , sortable: false}
      ],
    };
  },
  mounted() {
    this.fetchItems();
  },
  created() {
    setInterval(this.fetchItems, 30000);
  },
  methods: {
    fetchItems() {
      if(this.isLoading){
        return
      }
      this.isLoading = true

      ///axios.get(`http://127.0.0.1:8090/api/list-queue?withParentComplete=${this.withCompleteParent}&withChildrenComplete=true`)
      axios.get(`http://80.76.35.119/api/list-queue?withParentComplete=${this.withCompleteParent}&withChildrenComplete=true`)
        .then(response => {
          this.items = response.data.items;
        })
        .catch(error => {
          console.error('Ошибка при получении данных:', error);
        });

      this.isLoading = false
    },
    setIsComplete(item) {
      axios.post('/api/set-complete', { id: item.id, state:item.is_complete })
        .then(response => {
          this.fetchItems()
        })
        .catch(error => {
          console.error('Ошибка при обновлении флага:', error);
        });
    },
    setIsPrinting(item) {
      axios.post('/api/set-printing', { id: item.id, state:item.is_printing })
        .then(response => {})
        .catch(error => {
          console.error('Ошибка при обновлении флага:', error);
        });
    },
  },
};
</script>
