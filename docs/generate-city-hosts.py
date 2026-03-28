#!/usr/bin/env python3
"""生成城市二级域名 hosts 配置（优先全量 city_list，回退 seed）。"""

import json
import os
import re
import subprocess
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parents[1]
SEED_PATH = BASE_DIR / "backend" / "internal" / "bootstrap" / "data" / "city_list.seed.json"
ENV_PATH = BASE_DIR / "backend" / ".env"


def simple_pinyin(name: str) -> str:
    """
    极简拼音转换：只保留英文字母数字，其他替换为下划线，然后缩减连续下划线。
    对于当前种子数据（纯中文城市名），通常前端会有更专业的拼音库；
    这里为了生成 hosts，仅做非常粗略的“slug 化”，避免额外依赖。
    """
    # 占位：中文统一用拼音占位符，避免生成空串
    # 简单映射常见城市提升可读性
    mapping = {
        "北京市": "beijing",
        "天津市": "tianjin",
        "石家庄市": "shijiazhuang",
        "太原市": "taiyuan",
        "呼和浩特市": "huhehaote",
        "沈阳市": "shenyang",
        "长春市": "changchun",
        "哈尔滨市": "haerbin",
        "上海市": "shanghai",
        "南京市": "nanjing",
        "杭州市": "hangzhou",
        "合肥市": "hefei",
        "福州市": "fuzhou",
        "南昌市": "nanchang",
        "济南市": "jinan",
        "郑州市": "zhengzhou",
        "武汉市": "wuhan",
        "长沙市": "changsha",
        "广州市": "guangzhou",
        "深圳市": "shenzhen",
        "南宁市": "nanning",
        "海口市": "haikou",
        "重庆市": "chongqing",
        "成都市": "chengdu",
        "贵阳市": "guiyang",
        "昆明市": "kunming",
        "拉萨市": "lasa",
        "西安市": "xian",
        "兰州市": "lanzhou",
        "西宁市": "xining",
        "银川市": "yinchuan",
        "乌鲁木齐市": "wulumuqi",
        "大连市": "dalian",
        "苏州市": "suzhou",
        "宁波市": "ningbo",
        "厦门市": "xiamen",
        "青岛市": "qingdao",
        "临沂市": "linyi",
        "洛阳市": "luoyang",
        "宜昌市": "yichang",
        "常德市": "changde",
        "佛山市": "foshan",
        "东莞市": "dongguan",
        "中山市": "zhongshan",
        "桂林市": "guilin",
        "绵阳市": "mianyang",
        "南充市": "nanchong",
        "毕节市": "bijie",
        "曲靖市": "qujing",
        "咸阳市": "xianyang",
    }
    if name in mapping:
        return mapping[name]

    # 兜底：把非字母数字都替换成下划线，再做精简
    slug = re.sub(r"[^a-zA-Z0-9]+", "_", name)
    slug = re.sub(r"_+", "_", slug).strip("_")
    return slug.lower() or "city"


def load_env_file(path: Path) -> dict:
    env = {}
    if not path.exists():
        return env
    for line in path.read_text(encoding="utf-8").splitlines():
        raw = line.strip()
        if not raw or raw.startswith("#") or "=" not in raw:
            continue
        k, v = raw.split("=", 1)
        env[k.strip()] = v.strip().strip('"').strip("'")
    return env


def _mysql_query(cmd_base: list[str], sql: str) -> str:
    return subprocess.check_output(cmd_base + ["-N", "-e", sql], text=True, stderr=subprocess.DEVNULL)


def try_load_cities_from_db() -> list[dict]:
    env_file = load_env_file(ENV_PATH)
    host = os.environ.get("DB_HOST") or env_file.get("DB_HOST", "127.0.0.1")
    port = os.environ.get("DB_PORT") or env_file.get("DB_PORT", "3306")
    user = (
        os.environ.get("DB_USER")
        or os.environ.get("DB_USERNAME")
        or env_file.get("DB_USER")
        or env_file.get("DB_USERNAME", "root")
    )
    password = os.environ.get("DB_PASSWORD") or env_file.get("DB_PASSWORD", "")
    name = (
        os.environ.get("DB_NAME")
        or os.environ.get("DB_DATABASE")
        or env_file.get("DB_NAME")
        or env_file.get("DB_DATABASE", "")
    )
    if not name:
        return []

    cmd_base = [
        "mysql",
        f"-h{host}",
        f"-P{port}",
        f"-u{user}",
        f"-D{name}",
    ]

    # 1) 优先 city_list（已清洗过的城市表，含 pinyin）
    query = (
        "SELECT city_code, name, pinyin "
        "FROM city_list "
        "WHERE status = 1 AND name IS NOT NULL AND name <> '' "
        "ORDER BY city_code ASC;"
    )
    env = os.environ.copy()
    env["MYSQL_PWD"] = password
    try:
        out = subprocess.check_output(cmd_base + ["-N", "-e", query], text=True, env=env, stderr=subprocess.DEVNULL)
    except Exception:
        out = ""

    rows: list[dict] = []
    for line in out.splitlines():
        parts = line.split("\t")
        if len(parts) < 3:
            continue
        city_code, city_name, city_pinyin = parts[0].strip(), parts[1].strip(), parts[2].strip()
        if not city_name:
            continue
        rows.append(
            {
                "cityCode": int(city_code) if city_code.isdigit() else 0,
                "name": city_name,
                "pinyin": city_pinyin,
            }
        )
    if rows:
        return rows

    # 2) 回退 area_code（如果用户导入了全量行政区数据）
    # 动态探测字段名，兼容不同建表脚本。
    try:
        cols_out = subprocess.check_output(
            cmd_base + ["-N", "-e", "SHOW COLUMNS FROM area_code;"],
            text=True,
            env=env,
            stderr=subprocess.DEVNULL,
        )
    except Exception:
        return []

    cols = [line.split("\t")[0].strip().lower() for line in cols_out.splitlines() if line.strip()]
    if not cols:
        return []
    name_col = "name" if "name" in cols else ("city_name" if "city_name" in cols else "")
    pinyin_col = "pinyin" if "pinyin" in cols else ("city_pinyin" if "city_pinyin" in cols else "")
    code_col = "code" if "code" in cols else ("city_code" if "city_code" in cols else ("id" if "id" in cols else ""))
    if not name_col:
        return []

    select_cols = [c for c in [code_col, name_col, pinyin_col] if c]
    if not select_cols:
        return []
    area_sql = f"SELECT {', '.join(select_cols)} FROM area_code ORDER BY {code_col or name_col} ASC;"
    try:
        out = subprocess.check_output(
            cmd_base + ["-N", "-e", area_sql],
            text=True,
            env=env,
            stderr=subprocess.DEVNULL,
        )
    except Exception:
        return []

    rows = []
    for line in out.splitlines():
        parts = [p.strip() for p in line.split("\t")]
        if not parts:
            continue
        city_code = parts[0] if code_col else ""
        city_name = parts[1] if code_col else parts[0]
        city_pinyin = parts[2] if pinyin_col and len(parts) > 2 else ""
        if not city_name:
            continue
        rows.append(
            {
                "cityCode": int(city_code) if str(city_code).isdigit() else 0,
                "name": city_name,
                "pinyin": city_pinyin,
            }
        )
    return rows


def load_cities() -> tuple[list[dict], str]:
    db_rows = try_load_cities_from_db()
    if db_rows:
        return db_rows, "city_list(db)"
    if not SEED_PATH.exists():
        raise SystemExit(f"seed file not found: {SEED_PATH}")
    return json.loads(SEED_PATH.read_text(encoding="utf-8")), "city_list.seed.json"


def main():
    data, source = load_cities()

    domain_suffix = os.environ.get("CITY_HOST_SUFFIX", "1.com")
    ip = os.environ.get("CITY_HOST_IP", "127.0.0.1")

    lines = set()
    for item in data:
        name = item.get("name") or ""
        py = (item.get("pinyin") or "").strip().lower()
        if not py:
            py = simple_pinyin(name)
        host = f"{py}.{domain_suffix}"
        lines.add(f"{ip} {host}")

    print(f"# Generated from {source}")
    print(f"# IP={ip}, suffix={domain_suffix}")
    for line in sorted(lines):
        print(line)


if __name__ == "__main__":
    main()

