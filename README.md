# memov2

Markdown ファイルベースのメモ・タスク管理 CLI/TUI ツール。

データベースを使わず、ファイルシステム上の Markdown ファイルとして直接保存するため、任意のエディタで編集でき、Git によるバージョン管理にも適しています。

## 特徴

- **ファイルベース**: Markdown ファイルとして保存。YAML frontmatter でメタデータを管理
- **カテゴリ分類**: 階層的なカテゴリツリーでメモを整理
- **タイムスタンプ管理**: 日時ベースのファイル名で自動整理
- **TUI ブラウザ**: ターミナル内でメモの閲覧・検索が可能
- **週次レポート**: メモ・タスクの週次サマリーを自動生成

## インストール

### ビルド

```bash
git clone https://github.com/hirotoni/memov2.git
cd memov2
make build    # ./memov2 にバイナリを生成
make install  # $GOPATH/bin にインストール
```

Go 1.24.3 以上が必要です。

## 使い方

### メモ

```bash
# 新しいメモを作成
memov2 memos new "ミーティングノート"

# メモを TUI で閲覧・検索
memov2 memos browse

# メモ一覧を表示
memov2 memos list

# メモを開く
memov2 memos open "memos/work/20250114Mon150405_memo_notes.md"

# ローマ字対応の検索（fzf との連携用）
memov2 memos search "kaigi"

# マッチコンテキスト付きで検索（マッチ箇所の詳細を表示）
memov2 memos search --context "kaigi"

# メモのリネーム（タイトルとファイル名を変更）
memov2 memos rename "work/20250114Mon150405_memo_notes.md" "新しいタイトル"
memov2 memos rename "work/20250114Mon150405_memo_notes.md"  # 対話的に入力

# 週次レポートを生成
memov2 memos weekly

# メモのインデックスファイルを生成
memov2 memos index
```

### タスク

```bash
# 今日のタスクファイルを作成
memov2 todos new

# 既存のタスクファイルを初期化して再作成
memov2 todos new --truncate

# タスクの週次レポートを生成
memov2 todos weekly
```

### 設定

```bash
# 現在の設定を表示
memov2 config show

# 設定ファイルをエディタで編集
memov2 config edit
```

## 設定ファイル

初回起動時に `~/.config/memov2/config.toml` が自動生成されます。

```toml
base_dir = "~/.config/memov2/dailymemo/"   # ファイル保存先のベースディレクトリ
memos_foldername = "memos/"                 # メモ用サブディレクトリ
todos_foldername = "todos/"                 # タスク用サブディレクトリ
todos_daystoseek = 10                       # タスク引き継ぎの遡り日数
```

エディタは環境変数 `$EDITOR` で設定します。

## ファイル形式

### メモファイル

ファイル名: `YYYYMMDDDAY000000_memo_title.md`

```markdown
---
category: ["work", "projects"]
---

# ミーティングノート

## 議題1

内容...

## 議題2

内容...
```

### タスクファイル

ファイル名: `YYYYMMDDDAY_todos.md`

```markdown
# 20250214Fri

## meetings

## todos

- [ ] タスクA
- [ ] タスクB
  - [ ] サブタスク1
  - [ ] サブタスク2
- [x] 完了タスク

## wanttodos
```

## TUI (memov2 memos browse)

ブラウズモードとサーチモードの2つのモードがあり、`Tab` で切り替えられます。

### ブラウズモード

カテゴリ別のツリー表示でメモを閲覧できます。

| キー | 操作 |
|------|------|
| `j` / `k` | 上下移動 |
| `Ctrl+u` / `Ctrl+d` | 10行ジャンプ |
| `l` | ディレクトリ展開 / ファイルを開く |
| `h` | ディレクトリ折りたたみ |
| `>` / `<` | 全展開 / 全折りたたみ |
| `p` | プレビューパネル切り替え |
| `N` | 選択カテゴリに新規メモ作成 |
| `r` | メモのリネーム |
| `d` | メモの複製 |
| `D` | メモの削除（ゴミ箱へ移動） |
| `c` | カテゴリ変更 |
| `Tab` | サーチモードへ切り替え |
| `q` | 終了 |

### サーチモード

タイトル・カテゴリ・本文・見出しを横断的にファジー検索できます。ローマ字による日本語検索にも対応しています。

## fzf を使ったインタラクティブ検索

`memov2 memos search` は、ローマ字から日本語への変換（SKK 辞書ベース）に対応した検索コマンドです。fzf と組み合わせることで、TUI を起動せずにインタラクティブなメモ検索ができます。

### 基本（タイトルとパスのみ）

```bash
memov2 memos list | fzf \
  --disabled \
  --bind "change:reload:memov2 memos search {q}" \
| cut -f2 | xargs memov2 memos open
```

**動作の流れ:**

1. `memov2 memos list` でメモの一覧を初期表示
2. 文字を入力すると `memov2 memos search {q}` が呼ばれ、リアルタイムに絞り込み
3. ローマ字入力（例: `kaigi`）で日本語のメモ（例: 「会議」を含むメモ）もヒット
4. 選択したメモのパスを `cut -f2` で抽出し、`memov2 memos open` で開く

### マッチコンテキスト付き（本文のどこが一致したか表示）

`--context` (`-c`) フラグを使うと、各マッチの種類と一致箇所を表示できます。

```bash
memov2 memos list | fzf \
  --disabled \
  --delimiter=$'\t' \
  --with-nth='1,3,4' \
  --bind "change:reload:memov2 memos search --context {q}" \
| cut -f2 | xargs memov2 memos open
```

**出力フォーマット（タブ区切り4フィールド）:**

```
タイトル	path/to/file.md	[Content]	## 設計 > 一致した本文の行
タイトル	path/to/file.md	[Title]  	タイトル
タイトル	path/to/file.md	[Heading]	見出しテキスト
```

- `--with-nth='1,3,4'` でタイトル・マッチタイプ・マッチ内容を表示（パスは非表示）
- `cut -f2` でパスを抽出し `memos open` に渡す

### fzf でメモを選んでリネーム

```bash
memov2 memos list | fzf | cut -f2 | xargs memov2 memos rename
```

fzf でメモを選択すると、現在のタイトルが表示され、新しいタイトルを対話的に入力できます。

### シェル関数の例

```bash
# .bashrc / .zshrc に追加
memo-search() {
  memov2 memos list | fzf \
    --disabled \
    --bind "change:reload:memov2 memos search {q}" \
  | cut -f2 | xargs memov2 memos open
}

# マッチコンテキスト付き検索（本文の一致箇所を確認しながら選択）
memo-search-context() {
  memov2 memos list | fzf \
    --disabled \
    --delimiter=$'\t' \
    --with-nth='1,3,4' \
    --bind "change:reload:memov2 memos search --context {q}" \
  | cut -f2 | xargs memov2 memos open
}

memo-rename() {
  memov2 memos list | fzf | cut -f2 | xargs memov2 memos rename
}
```

## アーキテクチャ

```
cmd/                 → Cobra CLI コマンド定義
internal/
  app/               → DI コンテナ（設定・サービス・リポジトリの組み立て）
  service/           → ビジネスロジック（memo, todo, config）
  repositories/      → データアクセス（ファイルシステム操作）
  platform/          → 外部連携（エディタ, ファイルシステム, ゴミ箱）
  domain/            → エンティティ（MemoFile, TodoFile, WeeklyFile）
  interfaces/        → インターフェース定義（全レイヤー共通）
  ui/tui/            → Bubbletea TUI（ブラウズ + サーチ）
  config/            → TOML ベースの設定管理
  common/            → エラーハンドリング, ロギング
```

## 開発

```bash
make build    # ビルド
make install  # インストール
make test     # 全テスト実行（カバレッジ付き）
```

## ライセンス

MIT
