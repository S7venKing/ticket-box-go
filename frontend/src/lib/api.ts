const baseURL = import.meta.env.VITE_API_URL ?? "http://localhost:8081";

function messageForResponse(status: number, body: string): string {
  const normalized = body.trim().toLowerCase();
  if (normalized.includes("invalid credentials")) return "Email hoặc mật khẩu không đúng.";
  if (status === 401) return "Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.";
  if (status === 403 || normalized.includes("forbidden") || normalized.includes("admin role")) {
    return "Bạn không có quyền thực hiện thao tác này.";
  }
  if (status === 404 || normalized.includes("not found")) return "Không tìm thấy dữ liệu yêu cầu.";
  if (status === 405) return "Thao tác này không được hỗ trợ.";
  if (status === 400 && normalized.includes("invalid request")) return "Thông tin nhập chưa hợp lệ. Vui lòng kiểm tra lại.";
  if (status >= 500 || !body) return "Có lỗi xảy ra, vui lòng liên hệ quản trị viên.";
  return "Không thể hoàn tất thao tác. Vui lòng kiểm tra lại thông tin.";
}

export async function api<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem("access_token");
  let response: Response;
  try {
    response = await fetch(`${baseURL}${path}`, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...options.headers,
      },
    });
  } catch (error) {
    if (error instanceof TypeError) {
      // eslint-disable-next-line preserve-caught-error
      throw new Error("Không thể kết nối máy chủ. Vui lòng liên hệ quản trị viên.");
    }
    throw error;
  }
  if (!response.ok) {
    const message = await response.text();
    throw new Error(messageForResponse(response.status, message));
  }
  return response.status === 204 ? (undefined as T) : response.json();
}

export const get = <T>(path: string) => api<T>(path);
export const post = <T>(path: string, body?: unknown) =>
  api<T>(path, { method: "POST", body: body === undefined ? undefined : JSON.stringify(body) });
export const put = <T>(path: string, body: unknown) =>
  api<T>(path, { method: "PUT", body: JSON.stringify(body) });
