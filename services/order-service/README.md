my-project/
├── cmd/
│ └── api/
│ └── main.go # จุดเริ่มต้นของแอปพลิเคชัน (Entry point)
├── internal/ # โค้ดที่ไม่ต้องการให้ package ภายนอกนำไปใช้
│ ├── handlers/ # รับ HTTP Request และส่ง Response (Controller)
│ ├── models/ # GORM Structs (ตารางใน Database)
│ ├── repositories/ # การจัดการ Database ผ่าน GORM (Query logic)
│ ├── services/ # Business Logic (เชื่อมระหว่าง Handlers และ Repositories)
│ └── database/ # การเชื่อมต่อ DB และการทำ AutoMigration
├── pkg/ # Code ที่อนุญาตให้โปรเจกต์อื่นดึงไปใช้ได้ (เช่น utils)
├── configs/ # ไฟล์ตั้งค่าต่างๆ (.env, config.yaml)
├── go.mod
└── go.sum
