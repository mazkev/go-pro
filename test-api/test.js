// Konfigurasi API
const BASE_URL = "https://api.b.ai/v1"; // Ganti dengan Base URL API Anda (Otomatis ditambahkan https:// jika tidak ada)
const API_KEY = "sk-6l7ijv9ptkotnbk515bca7dcrk8bhfyv"; // Ganti dengan API Key Anda
const MODEL = ""; // KOSONGKAN (biarkan '') untuk otomatis memilih model pertama yang tersedia!

// Otomatis koreksi URL agar memiliki skema https:// jika dilewatkan
function getNormalizedUrl(url) {
  let cleaned = (url || "").trim().replace(/\/+$/, "");

  // Strip common endpoint paths
  cleaned = cleaned.replace(/\/chat\/completions\/?$/i, "");
  cleaned = cleaned.replace(/\/completions\/?$/i, "");
  cleaned = cleaned.replace(/\/messages\/?$/i, "");
  cleaned = cleaned.replace(/\/models\/?$/i, "");

  if (!/^https?:\/\//i.test(cleaned)) {
    cleaned = "https://" + cleaned;
  }

  // Jika URL tidak memiliki penanda path versi, tambahkan /v1 secara default
  try {
    const urlObj = new URL(cleaned);
    if (urlObj.pathname === "/" || urlObj.pathname === "") {
      cleaned = cleaned + "/v1";
    }
  } catch (e) {
    // Abaikan jika parsing gagal
  }

  return cleaned;
}

const NORMALIZED_BASE_URL = getNormalizedUrl(BASE_URL);

/**
 * Mengambil daftar model yang tersedia untuk mengecek validitas API Key
 * @returns {Promise<string[]>} Daftar ID model yang ditemukan
 */
async function getAvailableModels() {
  console.log("=== Menguji Koneksi & Mencari Daftar Model ===");
  try {
    const response = await fetch(`${NORMALIZED_BASE_URL}/models`, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${API_KEY}`,
        "Content-Type": "application/json",
      },
    });

    console.log(`Status HTTP: ${response.status} ${response.statusText}`);
    const data = await response.json();

    if (response.ok) {
      const models = data.data || [];
      if (models.length === 0) {
        console.log("⚠️ Koneksi sukses, tetapi tidak ada model yang terdaftar.");
        return [];
      }

      console.log("✅ Koneksi Sukses! Model yang ditemukan:");
      models.forEach((model) => {
        console.log(` - ID Model: ${model.id} (Owned by: ${model.owned_by || "unknown"})`);
      });

      // Mengembalikan array ID model
      return models.map((model) => model.id);
    } else {
      console.error("❌ Gagal terhubung:", data.error || data);
      return [];
    }
  } catch (error) {
    console.error("❌ Error Jaringan/Koneksi:", error.message);
    return [];
  }
}

/**
 * Mengirim pesan obrolan (chat completion) untuk mengetes respons model
 * @param {string} modelId ID model yang akan digunakan
 */
async function testChatCompletion(modelId) {
  console.log(`\n=== Menguji Chat Completion dengan Model: "${modelId}" ===`);
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 15000); // 15s timeout

  try {
    const payload = {
      model: modelId,
      messages: [{ role: "user", content: 'Katakan "Halo Dunia" dalam 3 kata saja.' }],
    };

    const response = await fetch(`${NORMALIZED_BASE_URL}/chat/completions`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${API_KEY}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
      signal: controller.signal,
    });

    clearTimeout(timeoutId);
    console.log(`Status HTTP: ${response.status} ${response.statusText}`);
    const data = await response.json();

    if (response.ok) {
      console.log("✅ Pengujian Chat Berhasil!");
      console.log("Respons Model:");
      console.log(`----------------------------------------`);
      console.log(data.choices[0].message.content);
      console.log(`----------------------------------------`);
      if (data.choices[0].message.reasoning) {
        console.log("Proses Penalaran (Reasoning):");
        console.log(data.choices[0].message.reasoning);
        console.log(`----------------------------------------`);
      }
    } else {
      console.error("❌ Gagal melakukan chat:", data.error || data);
    }
  } catch (error) {
    clearTimeout(timeoutId);
    if (error.name === "AbortError") {
      console.error("❌ Error: Permintaan chat menggantung (Timeout setelah 15 detik).");
      console.error("💡 Ini berarti saluran hulu (upstream channel) untuk model ini di server proxy sedang tidak merespon.");
    } else {
      console.error("❌ Error Jaringan/Koneksi:", error.message);
    }
  }
}

// Jalankan pengujian
async function run() {
  // 1. Ambil daftar model
  const availableModels = await getAvailableModels();

  if (availableModels.length === 0) {
    console.error("❌ Tidak bisa melanjutkan tes chat karena tidak ada model yang ditemukan.");
    return;
  }

  // 2. Tentukan model yang akan dipakai
  let modelToUse = MODEL;
  if (!modelToUse) {
    modelToUse = availableModels[0];
    console.log(`💡 Info: Variabel MODEL kosong. Otomatis menggunakan model pertama: "${modelToUse}"`);
  }

  // 3. Jalankan tes chat
  await testChatCompletion(modelToUse);
}

run();
