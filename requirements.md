# テックラジオ 要件定義

## コンセプト
技術記事を毎日自動収集し、Claudeが「本質的な記事だけ」を選別してラジオ台本を作る。
ユーザーが台本をTTSで読み上げてラジオ完成。気になった記事はClaudeに深掘り質問する。

## 全体アーキテクチャ

```
GitHub Actions (毎日 cron トリガー。PCの電源状態に無関係にクラウドで実行)
  │
  ├─ 1. 記事収集
  │     ソース: Zenn/Qiitaの go・aws・ai タグ + 各分野の公式ブログRSS
  │     (はてなブックマーク全体のような雑多なソースは含めない。関心領域に絞る)
  │
  ├─ 2. Claude API 呼び出し(claude-sonnet-5)
  │     - 各記事の本文を取得して読み込む
  │     - 選別基準: 単なる機能アップデート・製品リリースの告知は除外。
  │       設計判断・失敗談・実践知・組織論など「中身のある」記事を採用。
  │     - 台本を生成。記事本文の引用は避け、要約は自分の言葉でパラフレーズする(著作権配慮)
  │     - 採用/不採用の理由も構造化データとして出力する(後述の改善サイクル用)
  │
  ├─ 3. 保存
  │     - 台本: archive/YYYY-MM-DD.md としてリポジトリにコミット
  │     - メタデータ: archive/index.json に日付・採用記事リスト・不採用記事リスト・
  │       採用/不採用の理由を追記していく(蓄積してfew-shotや傾向分析に使う)
  │
  └─ 4. Discord Webhook で完成通知(Botではなく単純なWebhook。
        「今日の台本できたよ」+ archiveファイルのGitHub URLのみを送る)
```

## 深掘りフロー(Claude.ai Pro Project経由)

Discord BotでAPIを直接叩く方式は不採用。理由は、API経由の会話とClaude.ai(Pro)の会話は
別システムで記憶が連携しないため。代わりに以下のフローにする。

```
Discord通知(URLのみ) をiPhoneで見る
  → Claude.aiアプリの専用Project「tech radio」を開く
  → 「今日の台本読んで、○番目の記事を深掘りして」と話しかける
  → ClaudeがGitHubのarchiveファイルURLをfetchして読み込み、回答する
```

- API課金は発生しない(Pro/Maxの定額範囲内)
- PCの起動状態に依存しない(Remote Controlのような制約なし)
- Project Instructionsに以下を登録しておく:
  - archiveリポジトリのURLパターン(`https://github.com/<user>/tech-radio/blob/main/archive/{日付}.md`)
  - プロジェクトの背景・選別基準の要約
  - これにより「今日の台本」と言うだけでClaudeが日付からURLを組み立てて読みに行ける
- 会話はすべてこのProject内に蓄積されるため、過去の深掘りの続きも自然に繋がる
- Remote Control(iPhoneからPCのClaude Codeセッションに接続する公式機能)は、
  このプロジェクトの本筋には使わない。がっつりコード修正したい時のみのオプション

## 決定事項

| 項目 | 内容 |
|---|---|
| 実行基盤 | GitHub Actions(scheduled workflow)。Claude CodeのセッションスコープなcronやDesktopのスケジュールタスクはPC依存のため不採用 |
| 通知 | Discord Webhook(Botなし、URL通知のみ。LINE Notifyは2025年3月末で終了済みのため不採用) |
| モデル | Claude Sonnet 5(claude-sonnet-5) |
| 想定コスト | 台本生成: 1回あたり約$0.11、毎日実行で月$3〜4程度(500円以下)。深掘りはClaude.ai Pro定額内でAPI課金なし |
| 記事ソース | Zenn/Qiitaの go・aws・ai タグ、各分野の公式ブログRSS |
| 台本保存先 | GitHubリポジトリの archive/ ディレクトリにコミット(履歴として蓄積) |
| 深掘り手段 | Claude.aiアプリ(Pro)の専用Project「tech radio」経由。Discord Bot/API直叩きは不採用(記憶が分断されるため) |

## 改善サイクルの設計(重要)

蓄積したarchive/index.jsonのデータは、Anthropicのモデル学習には使われない
(API利用はデフォルトでモデル学習に不使用)。あくまで**このプロジェクト内でのClaude活用を
賢くする**ためのローカルなフィードバックループとして使う。

- 過去の採用/不採用記事を数件、選別プロンプトにfew-shotとして埋め込み、基準のブレを抑える
- 半年〜1年分溜まったら、関心の変遷(golang/aws/aiのどれに注目していたか)を分析できるようにする
- 台本の良し悪しを月次で振り返り、選別プロンプトやトーンを調整する

## 未実装・要検討事項
- 記事ソースの正確なRSS URL/タグ名の確定(実装時に検証しながら決める)
- GitHub Actions workflow.yml の具体的な記述
- 記事収集スクリプト(Python想定、feedparser等)
- Claude APIへのプロンプト設計(選別基準・台本フォーマット・出力JSON構造)
- Discord Webhook URLの取得・GitHub Secretsへの登録
- archive/index.json のスキーマ設計
- Claude.aiにProject「tech radio」を作成し、Project Instructionsに
  archiveリポジトリのURLパターンと選別基準を登録する
