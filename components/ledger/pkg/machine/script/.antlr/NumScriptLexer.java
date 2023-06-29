// Generated from /home/phlimy/Projects/Contrib/stack/components/ledger/pkg/machine/script/NumScript.g4 by ANTLR 4.9.2
import org.antlr.v4.runtime.Lexer;
import org.antlr.v4.runtime.CharStream;
import org.antlr.v4.runtime.Token;
import org.antlr.v4.runtime.TokenStream;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.misc.*;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast"})
public class NumScriptLexer extends Lexer {
	static { RuntimeMetaData.checkVersion("4.9.2", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		T__0=1, T__1=2, T__2=3, T__3=4, NEWLINE=5, WHITESPACE=6, MULTILINE_COMMENT=7, 
		LINE_COMMENT=8, VARS=9, META=10, SET_TX_META=11, SET_ACCOUNT_META=12, 
		PRINT=13, FAIL=14, SEND=15, SOURCE=16, FROM=17, MAX=18, DESTINATION=19, 
		TO=20, ALLOCATE=21, OP_ADD=22, OP_SUB=23, LPAREN=24, RPAREN=25, LBRACK=26, 
		RBRACK=27, LBRACE=28, RBRACE=29, EQ=30, TY_ACCOUNT=31, TY_ASSET=32, TY_NUMBER=33, 
		TY_MONETARY=34, TY_PORTION=35, TY_STRING=36, STRING=37, PORTION=38, REMAINING=39, 
		KEPT=40, BALANCE=41, SAVE=42, NUMBER=43, PERCENT=44, VARIABLE_NAME=45, 
		ACCOUNT=46, ASSET=47;
	public static String[] channelNames = {
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN"
	};

	public static String[] modeNames = {
		"DEFAULT_MODE"
	};

	private static String[] makeRuleNames() {
		return new String[] {
			"T__0", "T__1", "T__2", "T__3", "NEWLINE", "WHITESPACE", "MULTILINE_COMMENT", 
			"LINE_COMMENT", "VARS", "META", "SET_TX_META", "SET_ACCOUNT_META", "PRINT", 
			"FAIL", "SEND", "SOURCE", "FROM", "MAX", "DESTINATION", "TO", "ALLOCATE", 
			"OP_ADD", "OP_SUB", "LPAREN", "RPAREN", "LBRACK", "RBRACK", "LBRACE", 
			"RBRACE", "EQ", "TY_ACCOUNT", "TY_ASSET", "TY_NUMBER", "TY_MONETARY", 
			"TY_PORTION", "TY_STRING", "STRING", "PORTION", "REMAINING", "KEPT", 
			"BALANCE", "SAVE", "NUMBER", "PERCENT", "VARIABLE_NAME", "ACCOUNT", "ASSET"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, "'*'", "'allowing overdraft up to'", "'allowing unbounded overdraft'", 
			"','", null, null, null, null, "'vars'", "'meta'", "'set_tx_meta'", "'set_account_meta'", 
			"'print'", "'fail'", "'send'", "'source'", "'from'", "'max'", "'destination'", 
			"'to'", "'allocate'", "'+'", "'-'", "'('", "')'", "'['", "']'", "'{'", 
			"'}'", "'='", "'account'", "'asset'", "'number'", "'monetary'", "'portion'", 
			"'string'", null, null, "'remaining'", "'kept'", "'balance'", "'save'", 
			null, "'%'"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, null, null, null, null, "NEWLINE", "WHITESPACE", "MULTILINE_COMMENT", 
			"LINE_COMMENT", "VARS", "META", "SET_TX_META", "SET_ACCOUNT_META", "PRINT", 
			"FAIL", "SEND", "SOURCE", "FROM", "MAX", "DESTINATION", "TO", "ALLOCATE", 
			"OP_ADD", "OP_SUB", "LPAREN", "RPAREN", "LBRACK", "RBRACK", "LBRACE", 
			"RBRACE", "EQ", "TY_ACCOUNT", "TY_ASSET", "TY_NUMBER", "TY_MONETARY", 
			"TY_PORTION", "TY_STRING", "STRING", "PORTION", "REMAINING", "KEPT", 
			"BALANCE", "SAVE", "NUMBER", "PERCENT", "VARIABLE_NAME", "ACCOUNT", "ASSET"
		};
	}
	private static final String[] _SYMBOLIC_NAMES = makeSymbolicNames();
	public static final Vocabulary VOCABULARY = new VocabularyImpl(_LITERAL_NAMES, _SYMBOLIC_NAMES);

	/**
	 * @deprecated Use {@link #VOCABULARY} instead.
	 */
	@Deprecated
	public static final String[] tokenNames;
	static {
		tokenNames = new String[_SYMBOLIC_NAMES.length];
		for (int i = 0; i < tokenNames.length; i++) {
			tokenNames[i] = VOCABULARY.getLiteralName(i);
			if (tokenNames[i] == null) {
				tokenNames[i] = VOCABULARY.getSymbolicName(i);
			}

			if (tokenNames[i] == null) {
				tokenNames[i] = "<INVALID>";
			}
		}
	}

	@Override
	@Deprecated
	public String[] getTokenNames() {
		return tokenNames;
	}

	@Override

	public Vocabulary getVocabulary() {
		return VOCABULARY;
	}


	public NumScriptLexer(CharStream input) {
		super(input);
		_interp = new LexerATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	@Override
	public String getGrammarFileName() { return "NumScript.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public String[] getChannelNames() { return channelNames; }

	@Override
	public String[] getModeNames() { return modeNames; }

	@Override
	public ATN getATN() { return _ATN; }

	public static final String _serializedATN =
		"\3\u608b\ua72a\u8133\ub9ed\u417c\u3be7\u7786\u5964\2\61\u01cb\b\1\4\2"+
		"\t\2\4\3\t\3\4\4\t\4\4\5\t\5\4\6\t\6\4\7\t\7\4\b\t\b\4\t\t\t\4\n\t\n\4"+
		"\13\t\13\4\f\t\f\4\r\t\r\4\16\t\16\4\17\t\17\4\20\t\20\4\21\t\21\4\22"+
		"\t\22\4\23\t\23\4\24\t\24\4\25\t\25\4\26\t\26\4\27\t\27\4\30\t\30\4\31"+
		"\t\31\4\32\t\32\4\33\t\33\4\34\t\34\4\35\t\35\4\36\t\36\4\37\t\37\4 \t"+
		" \4!\t!\4\"\t\"\4#\t#\4$\t$\4%\t%\4&\t&\4\'\t\'\4(\t(\4)\t)\4*\t*\4+\t"+
		"+\4,\t,\4-\t-\4.\t.\4/\t/\4\60\t\60\3\2\3\2\3\3\3\3\3\3\3\3\3\3\3\3\3"+
		"\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3"+
		"\3\3\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3"+
		"\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\5\3\5\3\6\6\6\u009d"+
		"\n\6\r\6\16\6\u009e\3\7\6\7\u00a2\n\7\r\7\16\7\u00a3\3\7\3\7\3\b\3\b\3"+
		"\b\3\b\3\b\7\b\u00ad\n\b\f\b\16\b\u00b0\13\b\3\b\3\b\3\b\3\b\3\b\3\t\3"+
		"\t\3\t\3\t\7\t\u00bb\n\t\f\t\16\t\u00be\13\t\3\t\3\t\3\t\3\t\3\n\3\n\3"+
		"\n\3\n\3\n\3\13\3\13\3\13\3\13\3\13\3\f\3\f\3\f\3\f\3\f\3\f\3\f\3\f\3"+
		"\f\3\f\3\f\3\f\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r"+
		"\3\r\3\r\3\r\3\16\3\16\3\16\3\16\3\16\3\16\3\17\3\17\3\17\3\17\3\17\3"+
		"\20\3\20\3\20\3\20\3\20\3\21\3\21\3\21\3\21\3\21\3\21\3\21\3\22\3\22\3"+
		"\22\3\22\3\22\3\23\3\23\3\23\3\23\3\24\3\24\3\24\3\24\3\24\3\24\3\24\3"+
		"\24\3\24\3\24\3\24\3\24\3\25\3\25\3\25\3\26\3\26\3\26\3\26\3\26\3\26\3"+
		"\26\3\26\3\26\3\27\3\27\3\30\3\30\3\31\3\31\3\32\3\32\3\33\3\33\3\34\3"+
		"\34\3\35\3\35\3\36\3\36\3\37\3\37\3 \3 \3 \3 \3 \3 \3 \3 \3!\3!\3!\3!"+
		"\3!\3!\3\"\3\"\3\"\3\"\3\"\3\"\3\"\3#\3#\3#\3#\3#\3#\3#\3#\3#\3$\3$\3"+
		"$\3$\3$\3$\3$\3$\3%\3%\3%\3%\3%\3%\3%\3&\3&\7&\u0164\n&\f&\16&\u0167\13"+
		"&\3&\3&\3\'\6\'\u016c\n\'\r\'\16\'\u016d\3\'\5\'\u0171\n\'\3\'\3\'\5\'"+
		"\u0175\n\'\3\'\6\'\u0178\n\'\r\'\16\'\u0179\3\'\6\'\u017d\n\'\r\'\16\'"+
		"\u017e\3\'\3\'\6\'\u0183\n\'\r\'\16\'\u0184\5\'\u0187\n\'\3\'\5\'\u018a"+
		"\n\'\3(\3(\3(\3(\3(\3(\3(\3(\3(\3(\3)\3)\3)\3)\3)\3*\3*\3*\3*\3*\3*\3"+
		"*\3*\3+\3+\3+\3+\3+\3,\6,\u01a9\n,\r,\16,\u01aa\3-\3-\3.\3.\6.\u01b1\n"+
		".\r.\16.\u01b2\3.\7.\u01b6\n.\f.\16.\u01b9\13.\3/\3/\6/\u01bd\n/\r/\16"+
		"/\u01be\3/\7/\u01c2\n/\f/\16/\u01c5\13/\3\60\6\60\u01c8\n\60\r\60\16\60"+
		"\u01c9\4\u00ae\u00bc\2\61\3\3\5\4\7\5\t\6\13\7\r\b\17\t\21\n\23\13\25"+
		"\f\27\r\31\16\33\17\35\20\37\21!\22#\23%\24\'\25)\26+\27-\30/\31\61\32"+
		"\63\33\65\34\67\359\36;\37= ?!A\"C#E$G%I&K\'M(O)Q*S+U,W-Y.[/]\60_\61\3"+
		"\2\f\4\2\f\f\17\17\4\2\13\13\"\"\b\2\"\"//\62;C\\aac|\3\2\62;\3\2\"\""+
		"\4\2aac|\5\2\62;aac|\5\2C\\aac|\6\2\62<C\\aac|\4\2\61;C\\\2\u01de\2\3"+
		"\3\2\2\2\2\5\3\2\2\2\2\7\3\2\2\2\2\t\3\2\2\2\2\13\3\2\2\2\2\r\3\2\2\2"+
		"\2\17\3\2\2\2\2\21\3\2\2\2\2\23\3\2\2\2\2\25\3\2\2\2\2\27\3\2\2\2\2\31"+
		"\3\2\2\2\2\33\3\2\2\2\2\35\3\2\2\2\2\37\3\2\2\2\2!\3\2\2\2\2#\3\2\2\2"+
		"\2%\3\2\2\2\2\'\3\2\2\2\2)\3\2\2\2\2+\3\2\2\2\2-\3\2\2\2\2/\3\2\2\2\2"+
		"\61\3\2\2\2\2\63\3\2\2\2\2\65\3\2\2\2\2\67\3\2\2\2\29\3\2\2\2\2;\3\2\2"+
		"\2\2=\3\2\2\2\2?\3\2\2\2\2A\3\2\2\2\2C\3\2\2\2\2E\3\2\2\2\2G\3\2\2\2\2"+
		"I\3\2\2\2\2K\3\2\2\2\2M\3\2\2\2\2O\3\2\2\2\2Q\3\2\2\2\2S\3\2\2\2\2U\3"+
		"\2\2\2\2W\3\2\2\2\2Y\3\2\2\2\2[\3\2\2\2\2]\3\2\2\2\2_\3\2\2\2\3a\3\2\2"+
		"\2\5c\3\2\2\2\7|\3\2\2\2\t\u0099\3\2\2\2\13\u009c\3\2\2\2\r\u00a1\3\2"+
		"\2\2\17\u00a7\3\2\2\2\21\u00b6\3\2\2\2\23\u00c3\3\2\2\2\25\u00c8\3\2\2"+
		"\2\27\u00cd\3\2\2\2\31\u00d9\3\2\2\2\33\u00ea\3\2\2\2\35\u00f0\3\2\2\2"+
		"\37\u00f5\3\2\2\2!\u00fa\3\2\2\2#\u0101\3\2\2\2%\u0106\3\2\2\2\'\u010a"+
		"\3\2\2\2)\u0116\3\2\2\2+\u0119\3\2\2\2-\u0122\3\2\2\2/\u0124\3\2\2\2\61"+
		"\u0126\3\2\2\2\63\u0128\3\2\2\2\65\u012a\3\2\2\2\67\u012c\3\2\2\29\u012e"+
		"\3\2\2\2;\u0130\3\2\2\2=\u0132\3\2\2\2?\u0134\3\2\2\2A\u013c\3\2\2\2C"+
		"\u0142\3\2\2\2E\u0149\3\2\2\2G\u0152\3\2\2\2I\u015a\3\2\2\2K\u0161\3\2"+
		"\2\2M\u0189\3\2\2\2O\u018b\3\2\2\2Q\u0195\3\2\2\2S\u019a\3\2\2\2U\u01a2"+
		"\3\2\2\2W\u01a8\3\2\2\2Y\u01ac\3\2\2\2[\u01ae\3\2\2\2]\u01ba\3\2\2\2_"+
		"\u01c7\3\2\2\2ab\7,\2\2b\4\3\2\2\2cd\7c\2\2de\7n\2\2ef\7n\2\2fg\7q\2\2"+
		"gh\7y\2\2hi\7k\2\2ij\7p\2\2jk\7i\2\2kl\7\"\2\2lm\7q\2\2mn\7x\2\2no\7g"+
		"\2\2op\7t\2\2pq\7f\2\2qr\7t\2\2rs\7c\2\2st\7h\2\2tu\7v\2\2uv\7\"\2\2v"+
		"w\7w\2\2wx\7r\2\2xy\7\"\2\2yz\7v\2\2z{\7q\2\2{\6\3\2\2\2|}\7c\2\2}~\7"+
		"n\2\2~\177\7n\2\2\177\u0080\7q\2\2\u0080\u0081\7y\2\2\u0081\u0082\7k\2"+
		"\2\u0082\u0083\7p\2\2\u0083\u0084\7i\2\2\u0084\u0085\7\"\2\2\u0085\u0086"+
		"\7w\2\2\u0086\u0087\7p\2\2\u0087\u0088\7d\2\2\u0088\u0089\7q\2\2\u0089"+
		"\u008a\7w\2\2\u008a\u008b\7p\2\2\u008b\u008c\7f\2\2\u008c\u008d\7g\2\2"+
		"\u008d\u008e\7f\2\2\u008e\u008f\7\"\2\2\u008f\u0090\7q\2\2\u0090\u0091"+
		"\7x\2\2\u0091\u0092\7g\2\2\u0092\u0093\7t\2\2\u0093\u0094\7f\2\2\u0094"+
		"\u0095\7t\2\2\u0095\u0096\7c\2\2\u0096\u0097\7h\2\2\u0097\u0098\7v\2\2"+
		"\u0098\b\3\2\2\2\u0099\u009a\7.\2\2\u009a\n\3\2\2\2\u009b\u009d\t\2\2"+
		"\2\u009c\u009b\3\2\2\2\u009d\u009e\3\2\2\2\u009e\u009c\3\2\2\2\u009e\u009f"+
		"\3\2\2\2\u009f\f\3\2\2\2\u00a0\u00a2\t\3\2\2\u00a1\u00a0\3\2\2\2\u00a2"+
		"\u00a3\3\2\2\2\u00a3\u00a1\3\2\2\2\u00a3\u00a4\3\2\2\2\u00a4\u00a5\3\2"+
		"\2\2\u00a5\u00a6\b\7\2\2\u00a6\16\3\2\2\2\u00a7\u00a8\7\61\2\2\u00a8\u00a9"+
		"\7,\2\2\u00a9\u00ae\3\2\2\2\u00aa\u00ad\5\17\b\2\u00ab\u00ad\13\2\2\2"+
		"\u00ac\u00aa\3\2\2\2\u00ac\u00ab\3\2\2\2\u00ad\u00b0\3\2\2\2\u00ae\u00af"+
		"\3\2\2\2\u00ae\u00ac\3\2\2\2\u00af\u00b1\3\2\2\2\u00b0\u00ae\3\2\2\2\u00b1"+
		"\u00b2\7,\2\2\u00b2\u00b3\7\61\2\2\u00b3\u00b4\3\2\2\2\u00b4\u00b5\b\b"+
		"\2\2\u00b5\20\3\2\2\2\u00b6\u00b7\7\61\2\2\u00b7\u00b8\7\61\2\2\u00b8"+
		"\u00bc\3\2\2\2\u00b9\u00bb\13\2\2\2\u00ba\u00b9\3\2\2\2\u00bb\u00be\3"+
		"\2\2\2\u00bc\u00bd\3\2\2\2\u00bc\u00ba\3\2\2\2\u00bd\u00bf\3\2\2\2\u00be"+
		"\u00bc\3\2\2\2\u00bf\u00c0\5\13\6\2\u00c0\u00c1\3\2\2\2\u00c1\u00c2\b"+
		"\t\2\2\u00c2\22\3\2\2\2\u00c3\u00c4\7x\2\2\u00c4\u00c5\7c\2\2\u00c5\u00c6"+
		"\7t\2\2\u00c6\u00c7\7u\2\2\u00c7\24\3\2\2\2\u00c8\u00c9\7o\2\2\u00c9\u00ca"+
		"\7g\2\2\u00ca\u00cb\7v\2\2\u00cb\u00cc\7c\2\2\u00cc\26\3\2\2\2\u00cd\u00ce"+
		"\7u\2\2\u00ce\u00cf\7g\2\2\u00cf\u00d0\7v\2\2\u00d0\u00d1\7a\2\2\u00d1"+
		"\u00d2\7v\2\2\u00d2\u00d3\7z\2\2\u00d3\u00d4\7a\2\2\u00d4\u00d5\7o\2\2"+
		"\u00d5\u00d6\7g\2\2\u00d6\u00d7\7v\2\2\u00d7\u00d8\7c\2\2\u00d8\30\3\2"+
		"\2\2\u00d9\u00da\7u\2\2\u00da\u00db\7g\2\2\u00db\u00dc\7v\2\2\u00dc\u00dd"+
		"\7a\2\2\u00dd\u00de\7c\2\2\u00de\u00df\7e\2\2\u00df\u00e0\7e\2\2\u00e0"+
		"\u00e1\7q\2\2\u00e1\u00e2\7w\2\2\u00e2\u00e3\7p\2\2\u00e3\u00e4\7v\2\2"+
		"\u00e4\u00e5\7a\2\2\u00e5\u00e6\7o\2\2\u00e6\u00e7\7g\2\2\u00e7\u00e8"+
		"\7v\2\2\u00e8\u00e9\7c\2\2\u00e9\32\3\2\2\2\u00ea\u00eb\7r\2\2\u00eb\u00ec"+
		"\7t\2\2\u00ec\u00ed\7k\2\2\u00ed\u00ee\7p\2\2\u00ee\u00ef\7v\2\2\u00ef"+
		"\34\3\2\2\2\u00f0\u00f1\7h\2\2\u00f1\u00f2\7c\2\2\u00f2\u00f3\7k\2\2\u00f3"+
		"\u00f4\7n\2\2\u00f4\36\3\2\2\2\u00f5\u00f6\7u\2\2\u00f6\u00f7\7g\2\2\u00f7"+
		"\u00f8\7p\2\2\u00f8\u00f9\7f\2\2\u00f9 \3\2\2\2\u00fa\u00fb\7u\2\2\u00fb"+
		"\u00fc\7q\2\2\u00fc\u00fd\7w\2\2\u00fd\u00fe\7t\2\2\u00fe\u00ff\7e\2\2"+
		"\u00ff\u0100\7g\2\2\u0100\"\3\2\2\2\u0101\u0102\7h\2\2\u0102\u0103\7t"+
		"\2\2\u0103\u0104\7q\2\2\u0104\u0105\7o\2\2\u0105$\3\2\2\2\u0106\u0107"+
		"\7o\2\2\u0107\u0108\7c\2\2\u0108\u0109\7z\2\2\u0109&\3\2\2\2\u010a\u010b"+
		"\7f\2\2\u010b\u010c\7g\2\2\u010c\u010d\7u\2\2\u010d\u010e\7v\2\2\u010e"+
		"\u010f\7k\2\2\u010f\u0110\7p\2\2\u0110\u0111\7c\2\2\u0111\u0112\7v\2\2"+
		"\u0112\u0113\7k\2\2\u0113\u0114\7q\2\2\u0114\u0115\7p\2\2\u0115(\3\2\2"+
		"\2\u0116\u0117\7v\2\2\u0117\u0118\7q\2\2\u0118*\3\2\2\2\u0119\u011a\7"+
		"c\2\2\u011a\u011b\7n\2\2\u011b\u011c\7n\2\2\u011c\u011d\7q\2\2\u011d\u011e"+
		"\7e\2\2\u011e\u011f\7c\2\2\u011f\u0120\7v\2\2\u0120\u0121\7g\2\2\u0121"+
		",\3\2\2\2\u0122\u0123\7-\2\2\u0123.\3\2\2\2\u0124\u0125\7/\2\2\u0125\60"+
		"\3\2\2\2\u0126\u0127\7*\2\2\u0127\62\3\2\2\2\u0128\u0129\7+\2\2\u0129"+
		"\64\3\2\2\2\u012a\u012b\7]\2\2\u012b\66\3\2\2\2\u012c\u012d\7_\2\2\u012d"+
		"8\3\2\2\2\u012e\u012f\7}\2\2\u012f:\3\2\2\2\u0130\u0131\7\177\2\2\u0131"+
		"<\3\2\2\2\u0132\u0133\7?\2\2\u0133>\3\2\2\2\u0134\u0135\7c\2\2\u0135\u0136"+
		"\7e\2\2\u0136\u0137\7e\2\2\u0137\u0138\7q\2\2\u0138\u0139\7w\2\2\u0139"+
		"\u013a\7p\2\2\u013a\u013b\7v\2\2\u013b@\3\2\2\2\u013c\u013d\7c\2\2\u013d"+
		"\u013e\7u\2\2\u013e\u013f\7u\2\2\u013f\u0140\7g\2\2\u0140\u0141\7v\2\2"+
		"\u0141B\3\2\2\2\u0142\u0143\7p\2\2\u0143\u0144\7w\2\2\u0144\u0145\7o\2"+
		"\2\u0145\u0146\7d\2\2\u0146\u0147\7g\2\2\u0147\u0148\7t\2\2\u0148D\3\2"+
		"\2\2\u0149\u014a\7o\2\2\u014a\u014b\7q\2\2\u014b\u014c\7p\2\2\u014c\u014d"+
		"\7g\2\2\u014d\u014e\7v\2\2\u014e\u014f\7c\2\2\u014f\u0150\7t\2\2\u0150"+
		"\u0151\7{\2\2\u0151F\3\2\2\2\u0152\u0153\7r\2\2\u0153\u0154\7q\2\2\u0154"+
		"\u0155\7t\2\2\u0155\u0156\7v\2\2\u0156\u0157\7k\2\2\u0157\u0158\7q\2\2"+
		"\u0158\u0159\7p\2\2\u0159H\3\2\2\2\u015a\u015b\7u\2\2\u015b\u015c\7v\2"+
		"\2\u015c\u015d\7t\2\2\u015d\u015e\7k\2\2\u015e\u015f\7p\2\2\u015f\u0160"+
		"\7i\2\2\u0160J\3\2\2\2\u0161\u0165\7$\2\2\u0162\u0164\t\4\2\2\u0163\u0162"+
		"\3\2\2\2\u0164\u0167\3\2\2\2\u0165\u0163\3\2\2\2\u0165\u0166\3\2\2\2\u0166"+
		"\u0168\3\2\2\2\u0167\u0165\3\2\2\2\u0168\u0169\7$\2\2\u0169L\3\2\2\2\u016a"+
		"\u016c\t\5\2\2\u016b\u016a\3\2\2\2\u016c\u016d\3\2\2\2\u016d\u016b\3\2"+
		"\2\2\u016d\u016e\3\2\2\2\u016e\u0170\3\2\2\2\u016f\u0171\t\6\2\2\u0170"+
		"\u016f\3\2\2\2\u0170\u0171\3\2\2\2\u0171\u0172\3\2\2\2\u0172\u0174\7\61"+
		"\2\2\u0173\u0175\t\6\2\2\u0174\u0173\3\2\2\2\u0174\u0175\3\2\2\2\u0175"+
		"\u0177\3\2\2\2\u0176\u0178\t\5\2\2\u0177\u0176\3\2\2\2\u0178\u0179\3\2"+
		"\2\2\u0179\u0177\3\2\2\2\u0179\u017a\3\2\2\2\u017a\u018a\3\2\2\2\u017b"+
		"\u017d\t\5\2\2\u017c\u017b\3\2\2\2\u017d\u017e\3\2\2\2\u017e\u017c\3\2"+
		"\2\2\u017e\u017f\3\2\2\2\u017f\u0186\3\2\2\2\u0180\u0182\7\60\2\2\u0181"+
		"\u0183\t\5\2\2\u0182\u0181\3\2\2\2\u0183\u0184\3\2\2\2\u0184\u0182\3\2"+
		"\2\2\u0184\u0185\3\2\2\2\u0185\u0187\3\2\2\2\u0186\u0180\3\2\2\2\u0186"+
		"\u0187\3\2\2\2\u0187\u0188\3\2\2\2\u0188\u018a\7\'\2\2\u0189\u016b\3\2"+
		"\2\2\u0189\u017c\3\2\2\2\u018aN\3\2\2\2\u018b\u018c\7t\2\2\u018c\u018d"+
		"\7g\2\2\u018d\u018e\7o\2\2\u018e\u018f\7c\2\2\u018f\u0190\7k\2\2\u0190"+
		"\u0191\7p\2\2\u0191\u0192\7k\2\2\u0192\u0193\7p\2\2\u0193\u0194\7i\2\2"+
		"\u0194P\3\2\2\2\u0195\u0196\7m\2\2\u0196\u0197\7g\2\2\u0197\u0198\7r\2"+
		"\2\u0198\u0199\7v\2\2\u0199R\3\2\2\2\u019a\u019b\7d\2\2\u019b\u019c\7"+
		"c\2\2\u019c\u019d\7n\2\2\u019d\u019e\7c\2\2\u019e\u019f\7p\2\2\u019f\u01a0"+
		"\7e\2\2\u01a0\u01a1\7g\2\2\u01a1T\3\2\2\2\u01a2\u01a3\7u\2\2\u01a3\u01a4"+
		"\7c\2\2\u01a4\u01a5\7x\2\2\u01a5\u01a6\7g\2\2\u01a6V\3\2\2\2\u01a7\u01a9"+
		"\t\5\2\2\u01a8\u01a7\3\2\2\2\u01a9\u01aa\3\2\2\2\u01aa\u01a8\3\2\2\2\u01aa"+
		"\u01ab\3\2\2\2\u01abX\3\2\2\2\u01ac\u01ad\7\'\2\2\u01adZ\3\2\2\2\u01ae"+
		"\u01b0\7&\2\2\u01af\u01b1\t\7\2\2\u01b0\u01af\3\2\2\2\u01b1\u01b2\3\2"+
		"\2\2\u01b2\u01b0\3\2\2\2\u01b2\u01b3\3\2\2\2\u01b3\u01b7\3\2\2\2\u01b4"+
		"\u01b6\t\b\2\2\u01b5\u01b4\3\2\2\2\u01b6\u01b9\3\2\2\2\u01b7\u01b5\3\2"+
		"\2\2\u01b7\u01b8\3\2\2\2\u01b8\\\3\2\2\2\u01b9\u01b7\3\2\2\2\u01ba\u01bc"+
		"\7B\2\2\u01bb\u01bd\t\t\2\2\u01bc\u01bb\3\2\2\2\u01bd\u01be\3\2\2\2\u01be"+
		"\u01bc\3\2\2\2\u01be\u01bf\3\2\2\2\u01bf\u01c3\3\2\2\2\u01c0\u01c2\t\n"+
		"\2\2\u01c1\u01c0\3\2\2\2\u01c2\u01c5\3\2\2\2\u01c3\u01c1\3\2\2\2\u01c3"+
		"\u01c4\3\2\2\2\u01c4^\3\2\2\2\u01c5\u01c3\3\2\2\2\u01c6\u01c8\t\13\2\2"+
		"\u01c7\u01c6\3\2\2\2\u01c8\u01c9\3\2\2\2\u01c9\u01c7\3\2\2\2\u01c9\u01ca"+
		"\3\2\2\2\u01ca`\3\2\2\2\27\2\u009e\u00a3\u00ac\u00ae\u00bc\u0165\u016d"+
		"\u0170\u0174\u0179\u017e\u0184\u0186\u0189\u01aa\u01b2\u01b7\u01be\u01c3"+
		"\u01c9\3\b\2\2";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}