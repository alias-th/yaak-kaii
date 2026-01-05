# กำหนดอิมเมจพื้นฐาน (Base Image)
# alpine เป็นลินุกซ์ดิสโทรขนาดเล็กมาก
FROM alpine

# กำหนดไดเร็กทอรีทำงาน (Working Directory) ภายในคอนเทนเนอร์ (Container)
WORKDIR /app

# คัดลอกไฟล์/ไดเร็กทอรี จากเครื่องโฮสต์ (Host machine) ที่กำลังรันคำสั่ง docker build ไปยังอิมเมจ
ADD shared shared
ADD build build

# ดำเนินการรันไฟล์ที่ชื่อว่า api-gateway ซึ่งอยู่ในไดเร็กทอรี /app/build/
ENTRYPOINT build/order-service