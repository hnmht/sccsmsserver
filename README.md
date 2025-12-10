# Sea&Cloud Construction Site Management System
## How it all started
In 2015, I served as the CIO at a mining company. At that time, the company was located at a mine in central China, Where fatal safety accidents happened quite often. Behind every accident were broken families of the deceased, endless troubles for the company, and substantial financial losses.
I often visited the construction site at the mine, and I knew that root cause of these accidents lay in the "mindset of taking chances" shared by all participants. The construction areas were spread throuthout underground tunnels, making it impossilbe for managers to oversee every work face, Everyone hoped to do a little less work, believing that cutting corners probabley wouldn't lead to a safety incident. The solution is to leverage IT technology to upload photos of all pre-construction safety inspections to a system, making them publicly accessible for everyone to review. This ensures that every task can be supervised, effectively forcing everyone to perform their duties diligently.
When I brought up my idea to the boss, he didn't have high expectations, but he still approved the project, hoping to give it a try. We then spent three months developing a "Work Checklist" module on the company's ERP system. This module provided a browser-based interface that allowed employees to upload photos taken with digital cameras according to the pre-construction safety checklist. Everyone, including managers, could quickly view the checklist and related photos, and send email to the work site supervisor if any issues were identified. The results are very clear: after the system went live in October 2015, there were no fatal accidents at the throught 2016. This is the proudest project of my career.
## Why open source?
Three years after I left the mining company, from 2022 to 2023, I developed a "Field Management System". This system is a suite of application including server-side, browser-based, and mobile components, designed for private deployment. After nearly two years of promotion, I was disheartened to find that I couldn't sell the software in China. 
But I firmly believe this system has value. Perhaps somewhere in this vast world, someone needs it. Perhaps at some point in the furture, it could prevent a tragedy. That's why I've decided to open source it.

## Minio
minio:
  endpoint: "192.168.3.125:9000"
  accesskeyid: "minioadmin"
  secretaccesskey: "mht1977yeah"
  secure: false
  selfsigned: false
  defaultbucket: "seacloud"