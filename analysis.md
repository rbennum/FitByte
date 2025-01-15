# Load Test Analysis

This analysis was generated from ChatGPT after sending ours (monolith) and another team's setup (microservice).

__

| Metric                  | **DEBUG Team**                  | **Skibidi Team**                 | **Better Performer**      |
|-------------------------|---------------------------------------------|---------------------------------------------|---------------------------|
| **Checks**              | 88.88% (294,423/331,242)                   | 96.39% (2,771,194/2,874,814)                | **Skibidi Team**          |
| **Data Received**       | 65 MB (310 kB/s)                           | 342 MB (1.1 MB/s)                           | **Skibidi Team**          |
| **Data Sent**           | 380 MB (1.8 MB/s)                          | 2.2 GB (7.2 MB/s)                           | **Skibidi Team**          |
| **HTTP Req Duration**   | Avg: 2s, 90th: 5.07s, 95th: 8.31s          | Avg: 880ms, 90th: 2.21s, 95th: 5.46s        | **Skibidi Team**          |
| **HTTP Req Failed**     | 22.09% (23,330/105,569)                    | 11.57% (99,198/857,211)                     | **Skibidi Team**          |
| **HTTP Req Waiting**    | Avg: 2s, Min: 239.56ms, Max: 51.87s        | Avg: 880ms, Min: 0ms, Max: 1m0s             | **Skibidi Team**          |
| **HTTP Reqs (RPS)**     | 504.57 RPS                                 | 2,824 RPS                                   | **Skibidi Team**          |
| **Iteration Duration**  | Avg: 14.19s, 90th: 46.83s, 95th: 1m0s      | Avg: 15.42s, 90th: 48.73s, 95th: 1m1s       | **DEBUG Team (slightly)** |
| **Iterations**          | 6,709                                      | 49,085                                      | **Skibidi Team**          |
| **Virtual Users (VUs)** | Max: 6,000, Current: 4,905                 | Max: 6,000, Current: 626                    | **Skibidi Team (efficiency)** |
| **CPU Credits Used**    | **1 CPU Credit**                           | **10 CPU Credits**                          | **DEBUG Team (efficiency)** |
| **Requests per CPU Credit (R/C)** | **504.57 RPS / Credit**           | **282.4 RPS / Credit**                      | **DEBUG Team**            |

---

## **Analysis**

1. **Checks (Pass Rate):**
   - **Skibidi Team** has a much higher pass rate (96.39%) compared to the DEBUG Team (88.88%), indicating better reliability and fewer errors.

2. **Data Transfer:**
   - **Skibidi Team** shows significantly higher data transfer rates (1.1 MB/s received, 7.2 MB/s sent) compared to the DEBUG Team. This suggests it handled more requests successfully.

3. **HTTP Request Duration:**
   - Average request duration for the **Skibidi Team** is significantly lower (880ms) compared to the first (2s). Additionally, the 90th and 95th percentiles also show better performance.

4. **HTTP Failures:**
   - The **Skibidi Team** has a lower failure rate (11.57%) than the first (22.09%). This highlights better stability and fewer errors under load.

5. **Request Handling (RPS):**
   - **Skibidi Team** processed 2,824 requests per second (RPS) compared to the DEBUG Team’s 504 RPS, a significant improvement.

6. **Iteration Duration:**
   - The **Skibidi Team** has slightly faster iterations (14.19s avg) compared to the second (15.42s avg), but the difference is marginal. This is not a significant factor here.

7. **Virtual Users:**
   - Both tests had a maximum of 6,000 VUs, but the **Skibidi Team** had far fewer active VUs at its peak (626 compared to 4,905 in the first image). This suggests the second scenario is more efficient in resource usage while still achieving better performance.

8. **Requests per CPU Credit (R/C)**

   - This metric measures how many requests per second (RPS) were handled per CPU credit, highlighting CPU efficiency:

        - **DEBUG Team:** 504.57 RPS / 1 CPU Credit = **504.57 Requests/Credit**
        - **Skibidi Team:** 2,824 RPS / 10 CPU Credits = **282.4 Requests/Credit**

        Despite **Skibidi Team's** higher absolute performance, the **DEBUG Team** was **far more CPU-efficient**, handling nearly **1.78x more requests per CPU credit.**

---

## **Conclusion**

- **Overall Winner: Skibidi Team (Absolute Performance)**
  - Skibidi Team is significantly better in almost all key metrics: higher reliability (pass rate), better response times, higher throughput, and lower failure rates. It also handled significantly more requests per second with fewer active virtual users, indicating better scalability and efficiency.
  
- **CPU Efficiency Winner: DEBUG Team**
  - DEBUG Team achieved better CPU efficiency by processing significantly more requests per CPU credit. This is critical in environments where CPU credits are limited or expensive.

## **Recommendation:**

If CPU resources are abundant and cost is not a concern, **Skibidi Team's configuration** is preferable for maximum throughput. However, if efficiency is paramount, **DEBUG Team's setup** demonstrates excellent optimization for limited resources.
