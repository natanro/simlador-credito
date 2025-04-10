// init.js
db = db.getSiblingDB('simlador_credito'); // Substitua pelo nome do seu banco de dados

db.createCollection('parametros_classe_usuarios'); // Cria uma coleção chamada 'exemplo'

db.exemplo.insertMany([
    { classe: "0-25", taxa: 0.05 },
    { classe: "26-40", taxa: 0.03 },
    { classe: "41-60", taxa: 0.02 },
    { classe: "61+", taxa: 0.04 }
]);
