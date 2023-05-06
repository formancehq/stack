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
		KEPT=40, BALANCE=41, NUMBER=42, PERCENT=43, VARIABLE_NAME=44, ACCOUNT=45, 
		ASSET=46;
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
			"BALANCE", "NUMBER", "PERCENT", "VARIABLE_NAME", "ACCOUNT", "ASSET"
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
			"'string'", null, null, "'remaining'", "'kept'", "'balance'", null, "'%'"
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
			"BALANCE", "NUMBER", "PERCENT", "VARIABLE_NAME", "ACCOUNT", "ASSET"
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
		"\3\u608b\ua72a\u8133\ub9ed\u417c\u3be7\u7786\u5964\2\60\u01c4\b\1\4\2"+
		"\t\2\4\3\t\3\4\4\t\4\4\5\t\5\4\6\t\6\4\7\t\7\4\b\t\b\4\t\t\t\4\n\t\n\4"+
		"\13\t\13\4\f\t\f\4\r\t\r\4\16\t\16\4\17\t\17\4\20\t\20\4\21\t\21\4\22"+
		"\t\22\4\23\t\23\4\24\t\24\4\25\t\25\4\26\t\26\4\27\t\27\4\30\t\30\4\31"+
		"\t\31\4\32\t\32\4\33\t\33\4\34\t\34\4\35\t\35\4\36\t\36\4\37\t\37\4 \t"+
		" \4!\t!\4\"\t\"\4#\t#\4$\t$\4%\t%\4&\t&\4\'\t\'\4(\t(\4)\t)\4*\t*\4+\t"+
		"+\4,\t,\4-\t-\4.\t.\4/\t/\3\2\3\2\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3"+
		"\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\3\4\3"+
		"\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4"+
		"\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\4\3\5\3\5\3\6\6\6\u009b\n\6\r\6"+
		"\16\6\u009c\3\7\6\7\u00a0\n\7\r\7\16\7\u00a1\3\7\3\7\3\b\3\b\3\b\3\b\3"+
		"\b\7\b\u00ab\n\b\f\b\16\b\u00ae\13\b\3\b\3\b\3\b\3\b\3\b\3\t\3\t\3\t\3"+
		"\t\7\t\u00b9\n\t\f\t\16\t\u00bc\13\t\3\t\3\t\3\t\3\t\3\n\3\n\3\n\3\n\3"+
		"\n\3\13\3\13\3\13\3\13\3\13\3\f\3\f\3\f\3\f\3\f\3\f\3\f\3\f\3\f\3\f\3"+
		"\f\3\f\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r\3\r"+
		"\3\r\3\16\3\16\3\16\3\16\3\16\3\16\3\17\3\17\3\17\3\17\3\17\3\20\3\20"+
		"\3\20\3\20\3\20\3\21\3\21\3\21\3\21\3\21\3\21\3\21\3\22\3\22\3\22\3\22"+
		"\3\22\3\23\3\23\3\23\3\23\3\24\3\24\3\24\3\24\3\24\3\24\3\24\3\24\3\24"+
		"\3\24\3\24\3\24\3\25\3\25\3\25\3\26\3\26\3\26\3\26\3\26\3\26\3\26\3\26"+
		"\3\26\3\27\3\27\3\30\3\30\3\31\3\31\3\32\3\32\3\33\3\33\3\34\3\34\3\35"+
		"\3\35\3\36\3\36\3\37\3\37\3 \3 \3 \3 \3 \3 \3 \3 \3!\3!\3!\3!\3!\3!\3"+
		"\"\3\"\3\"\3\"\3\"\3\"\3\"\3#\3#\3#\3#\3#\3#\3#\3#\3#\3$\3$\3$\3$\3$\3"+
		"$\3$\3$\3%\3%\3%\3%\3%\3%\3%\3&\3&\7&\u0162\n&\f&\16&\u0165\13&\3&\3&"+
		"\3\'\6\'\u016a\n\'\r\'\16\'\u016b\3\'\5\'\u016f\n\'\3\'\3\'\5\'\u0173"+
		"\n\'\3\'\6\'\u0176\n\'\r\'\16\'\u0177\3\'\6\'\u017b\n\'\r\'\16\'\u017c"+
		"\3\'\3\'\6\'\u0181\n\'\r\'\16\'\u0182\5\'\u0185\n\'\3\'\5\'\u0188\n\'"+
		"\3(\3(\3(\3(\3(\3(\3(\3(\3(\3(\3)\3)\3)\3)\3)\3*\3*\3*\3*\3*\3*\3*\3*"+
		"\3+\6+\u01a2\n+\r+\16+\u01a3\3,\3,\3-\3-\6-\u01aa\n-\r-\16-\u01ab\3-\7"+
		"-\u01af\n-\f-\16-\u01b2\13-\3.\3.\6.\u01b6\n.\r.\16.\u01b7\3.\7.\u01bb"+
		"\n.\f.\16.\u01be\13.\3/\6/\u01c1\n/\r/\16/\u01c2\4\u00ac\u00ba\2\60\3"+
		"\3\5\4\7\5\t\6\13\7\r\b\17\t\21\n\23\13\25\f\27\r\31\16\33\17\35\20\37"+
		"\21!\22#\23%\24\'\25)\26+\27-\30/\31\61\32\63\33\65\34\67\359\36;\37="+
		" ?!A\"C#E$G%I&K\'M(O)Q*S+U,W-Y.[/]\60\3\2\f\4\2\f\f\17\17\4\2\13\13\""+
		"\"\b\2\"\"//\62;C\\aac|\3\2\62;\3\2\"\"\4\2aac|\5\2\62;aac|\5\2C\\aac"+
		"|\6\2\62<C\\aac|\4\2\61;C\\\2\u01d7\2\3\3\2\2\2\2\5\3\2\2\2\2\7\3\2\2"+
		"\2\2\t\3\2\2\2\2\13\3\2\2\2\2\r\3\2\2\2\2\17\3\2\2\2\2\21\3\2\2\2\2\23"+
		"\3\2\2\2\2\25\3\2\2\2\2\27\3\2\2\2\2\31\3\2\2\2\2\33\3\2\2\2\2\35\3\2"+
		"\2\2\2\37\3\2\2\2\2!\3\2\2\2\2#\3\2\2\2\2%\3\2\2\2\2\'\3\2\2\2\2)\3\2"+
		"\2\2\2+\3\2\2\2\2-\3\2\2\2\2/\3\2\2\2\2\61\3\2\2\2\2\63\3\2\2\2\2\65\3"+
		"\2\2\2\2\67\3\2\2\2\29\3\2\2\2\2;\3\2\2\2\2=\3\2\2\2\2?\3\2\2\2\2A\3\2"+
		"\2\2\2C\3\2\2\2\2E\3\2\2\2\2G\3\2\2\2\2I\3\2\2\2\2K\3\2\2\2\2M\3\2\2\2"+
		"\2O\3\2\2\2\2Q\3\2\2\2\2S\3\2\2\2\2U\3\2\2\2\2W\3\2\2\2\2Y\3\2\2\2\2["+
		"\3\2\2\2\2]\3\2\2\2\3_\3\2\2\2\5a\3\2\2\2\7z\3\2\2\2\t\u0097\3\2\2\2\13"+
		"\u009a\3\2\2\2\r\u009f\3\2\2\2\17\u00a5\3\2\2\2\21\u00b4\3\2\2\2\23\u00c1"+
		"\3\2\2\2\25\u00c6\3\2\2\2\27\u00cb\3\2\2\2\31\u00d7\3\2\2\2\33\u00e8\3"+
		"\2\2\2\35\u00ee\3\2\2\2\37\u00f3\3\2\2\2!\u00f8\3\2\2\2#\u00ff\3\2\2\2"+
		"%\u0104\3\2\2\2\'\u0108\3\2\2\2)\u0114\3\2\2\2+\u0117\3\2\2\2-\u0120\3"+
		"\2\2\2/\u0122\3\2\2\2\61\u0124\3\2\2\2\63\u0126\3\2\2\2\65\u0128\3\2\2"+
		"\2\67\u012a\3\2\2\29\u012c\3\2\2\2;\u012e\3\2\2\2=\u0130\3\2\2\2?\u0132"+
		"\3\2\2\2A\u013a\3\2\2\2C\u0140\3\2\2\2E\u0147\3\2\2\2G\u0150\3\2\2\2I"+
		"\u0158\3\2\2\2K\u015f\3\2\2\2M\u0187\3\2\2\2O\u0189\3\2\2\2Q\u0193\3\2"+
		"\2\2S\u0198\3\2\2\2U\u01a1\3\2\2\2W\u01a5\3\2\2\2Y\u01a7\3\2\2\2[\u01b3"+
		"\3\2\2\2]\u01c0\3\2\2\2_`\7,\2\2`\4\3\2\2\2ab\7c\2\2bc\7n\2\2cd\7n\2\2"+
		"de\7q\2\2ef\7y\2\2fg\7k\2\2gh\7p\2\2hi\7i\2\2ij\7\"\2\2jk\7q\2\2kl\7x"+
		"\2\2lm\7g\2\2mn\7t\2\2no\7f\2\2op\7t\2\2pq\7c\2\2qr\7h\2\2rs\7v\2\2st"+
		"\7\"\2\2tu\7w\2\2uv\7r\2\2vw\7\"\2\2wx\7v\2\2xy\7q\2\2y\6\3\2\2\2z{\7"+
		"c\2\2{|\7n\2\2|}\7n\2\2}~\7q\2\2~\177\7y\2\2\177\u0080\7k\2\2\u0080\u0081"+
		"\7p\2\2\u0081\u0082\7i\2\2\u0082\u0083\7\"\2\2\u0083\u0084\7w\2\2\u0084"+
		"\u0085\7p\2\2\u0085\u0086\7d\2\2\u0086\u0087\7q\2\2\u0087\u0088\7w\2\2"+
		"\u0088\u0089\7p\2\2\u0089\u008a\7f\2\2\u008a\u008b\7g\2\2\u008b\u008c"+
		"\7f\2\2\u008c\u008d\7\"\2\2\u008d\u008e\7q\2\2\u008e\u008f\7x\2\2\u008f"+
		"\u0090\7g\2\2\u0090\u0091\7t\2\2\u0091\u0092\7f\2\2\u0092\u0093\7t\2\2"+
		"\u0093\u0094\7c\2\2\u0094\u0095\7h\2\2\u0095\u0096\7v\2\2\u0096\b\3\2"+
		"\2\2\u0097\u0098\7.\2\2\u0098\n\3\2\2\2\u0099\u009b\t\2\2\2\u009a\u0099"+
		"\3\2\2\2\u009b\u009c\3\2\2\2\u009c\u009a\3\2\2\2\u009c\u009d\3\2\2\2\u009d"+
		"\f\3\2\2\2\u009e\u00a0\t\3\2\2\u009f\u009e\3\2\2\2\u00a0\u00a1\3\2\2\2"+
		"\u00a1\u009f\3\2\2\2\u00a1\u00a2\3\2\2\2\u00a2\u00a3\3\2\2\2\u00a3\u00a4"+
		"\b\7\2\2\u00a4\16\3\2\2\2\u00a5\u00a6\7\61\2\2\u00a6\u00a7\7,\2\2\u00a7"+
		"\u00ac\3\2\2\2\u00a8\u00ab\5\17\b\2\u00a9\u00ab\13\2\2\2\u00aa\u00a8\3"+
		"\2\2\2\u00aa\u00a9\3\2\2\2\u00ab\u00ae\3\2\2\2\u00ac\u00ad\3\2\2\2\u00ac"+
		"\u00aa\3\2\2\2\u00ad\u00af\3\2\2\2\u00ae\u00ac\3\2\2\2\u00af\u00b0\7,"+
		"\2\2\u00b0\u00b1\7\61\2\2\u00b1\u00b2\3\2\2\2\u00b2\u00b3\b\b\2\2\u00b3"+
		"\20\3\2\2\2\u00b4\u00b5\7\61\2\2\u00b5\u00b6\7\61\2\2\u00b6\u00ba\3\2"+
		"\2\2\u00b7\u00b9\13\2\2\2\u00b8\u00b7\3\2\2\2\u00b9\u00bc\3\2\2\2\u00ba"+
		"\u00bb\3\2\2\2\u00ba\u00b8\3\2\2\2\u00bb\u00bd\3\2\2\2\u00bc\u00ba\3\2"+
		"\2\2\u00bd\u00be\5\13\6\2\u00be\u00bf\3\2\2\2\u00bf\u00c0\b\t\2\2\u00c0"+
		"\22\3\2\2\2\u00c1\u00c2\7x\2\2\u00c2\u00c3\7c\2\2\u00c3\u00c4\7t\2\2\u00c4"+
		"\u00c5\7u\2\2\u00c5\24\3\2\2\2\u00c6\u00c7\7o\2\2\u00c7\u00c8\7g\2\2\u00c8"+
		"\u00c9\7v\2\2\u00c9\u00ca\7c\2\2\u00ca\26\3\2\2\2\u00cb\u00cc\7u\2\2\u00cc"+
		"\u00cd\7g\2\2\u00cd\u00ce\7v\2\2\u00ce\u00cf\7a\2\2\u00cf\u00d0\7v\2\2"+
		"\u00d0\u00d1\7z\2\2\u00d1\u00d2\7a\2\2\u00d2\u00d3\7o\2\2\u00d3\u00d4"+
		"\7g\2\2\u00d4\u00d5\7v\2\2\u00d5\u00d6\7c\2\2\u00d6\30\3\2\2\2\u00d7\u00d8"+
		"\7u\2\2\u00d8\u00d9\7g\2\2\u00d9\u00da\7v\2\2\u00da\u00db\7a\2\2\u00db"+
		"\u00dc\7c\2\2\u00dc\u00dd\7e\2\2\u00dd\u00de\7e\2\2\u00de\u00df\7q\2\2"+
		"\u00df\u00e0\7w\2\2\u00e0\u00e1\7p\2\2\u00e1\u00e2\7v\2\2\u00e2\u00e3"+
		"\7a\2\2\u00e3\u00e4\7o\2\2\u00e4\u00e5\7g\2\2\u00e5\u00e6\7v\2\2\u00e6"+
		"\u00e7\7c\2\2\u00e7\32\3\2\2\2\u00e8\u00e9\7r\2\2\u00e9\u00ea\7t\2\2\u00ea"+
		"\u00eb\7k\2\2\u00eb\u00ec\7p\2\2\u00ec\u00ed\7v\2\2\u00ed\34\3\2\2\2\u00ee"+
		"\u00ef\7h\2\2\u00ef\u00f0\7c\2\2\u00f0\u00f1\7k\2\2\u00f1\u00f2\7n\2\2"+
		"\u00f2\36\3\2\2\2\u00f3\u00f4\7u\2\2\u00f4\u00f5\7g\2\2\u00f5\u00f6\7"+
		"p\2\2\u00f6\u00f7\7f\2\2\u00f7 \3\2\2\2\u00f8\u00f9\7u\2\2\u00f9\u00fa"+
		"\7q\2\2\u00fa\u00fb\7w\2\2\u00fb\u00fc\7t\2\2\u00fc\u00fd\7e\2\2\u00fd"+
		"\u00fe\7g\2\2\u00fe\"\3\2\2\2\u00ff\u0100\7h\2\2\u0100\u0101\7t\2\2\u0101"+
		"\u0102\7q\2\2\u0102\u0103\7o\2\2\u0103$\3\2\2\2\u0104\u0105\7o\2\2\u0105"+
		"\u0106\7c\2\2\u0106\u0107\7z\2\2\u0107&\3\2\2\2\u0108\u0109\7f\2\2\u0109"+
		"\u010a\7g\2\2\u010a\u010b\7u\2\2\u010b\u010c\7v\2\2\u010c\u010d\7k\2\2"+
		"\u010d\u010e\7p\2\2\u010e\u010f\7c\2\2\u010f\u0110\7v\2\2\u0110\u0111"+
		"\7k\2\2\u0111\u0112\7q\2\2\u0112\u0113\7p\2\2\u0113(\3\2\2\2\u0114\u0115"+
		"\7v\2\2\u0115\u0116\7q\2\2\u0116*\3\2\2\2\u0117\u0118\7c\2\2\u0118\u0119"+
		"\7n\2\2\u0119\u011a\7n\2\2\u011a\u011b\7q\2\2\u011b\u011c\7e\2\2\u011c"+
		"\u011d\7c\2\2\u011d\u011e\7v\2\2\u011e\u011f\7g\2\2\u011f,\3\2\2\2\u0120"+
		"\u0121\7-\2\2\u0121.\3\2\2\2\u0122\u0123\7/\2\2\u0123\60\3\2\2\2\u0124"+
		"\u0125\7*\2\2\u0125\62\3\2\2\2\u0126\u0127\7+\2\2\u0127\64\3\2\2\2\u0128"+
		"\u0129\7]\2\2\u0129\66\3\2\2\2\u012a\u012b\7_\2\2\u012b8\3\2\2\2\u012c"+
		"\u012d\7}\2\2\u012d:\3\2\2\2\u012e\u012f\7\177\2\2\u012f<\3\2\2\2\u0130"+
		"\u0131\7?\2\2\u0131>\3\2\2\2\u0132\u0133\7c\2\2\u0133\u0134\7e\2\2\u0134"+
		"\u0135\7e\2\2\u0135\u0136\7q\2\2\u0136\u0137\7w\2\2\u0137\u0138\7p\2\2"+
		"\u0138\u0139\7v\2\2\u0139@\3\2\2\2\u013a\u013b\7c\2\2\u013b\u013c\7u\2"+
		"\2\u013c\u013d\7u\2\2\u013d\u013e\7g\2\2\u013e\u013f\7v\2\2\u013fB\3\2"+
		"\2\2\u0140\u0141\7p\2\2\u0141\u0142\7w\2\2\u0142\u0143\7o\2\2\u0143\u0144"+
		"\7d\2\2\u0144\u0145\7g\2\2\u0145\u0146\7t\2\2\u0146D\3\2\2\2\u0147\u0148"+
		"\7o\2\2\u0148\u0149\7q\2\2\u0149\u014a\7p\2\2\u014a\u014b\7g\2\2\u014b"+
		"\u014c\7v\2\2\u014c\u014d\7c\2\2\u014d\u014e\7t\2\2\u014e\u014f\7{\2\2"+
		"\u014fF\3\2\2\2\u0150\u0151\7r\2\2\u0151\u0152\7q\2\2\u0152\u0153\7t\2"+
		"\2\u0153\u0154\7v\2\2\u0154\u0155\7k\2\2\u0155\u0156\7q\2\2\u0156\u0157"+
		"\7p\2\2\u0157H\3\2\2\2\u0158\u0159\7u\2\2\u0159\u015a\7v\2\2\u015a\u015b"+
		"\7t\2\2\u015b\u015c\7k\2\2\u015c\u015d\7p\2\2\u015d\u015e\7i\2\2\u015e"+
		"J\3\2\2\2\u015f\u0163\7$\2\2\u0160\u0162\t\4\2\2\u0161\u0160\3\2\2\2\u0162"+
		"\u0165\3\2\2\2\u0163\u0161\3\2\2\2\u0163\u0164\3\2\2\2\u0164\u0166\3\2"+
		"\2\2\u0165\u0163\3\2\2\2\u0166\u0167\7$\2\2\u0167L\3\2\2\2\u0168\u016a"+
		"\t\5\2\2\u0169\u0168\3\2\2\2\u016a\u016b\3\2\2\2\u016b\u0169\3\2\2\2\u016b"+
		"\u016c\3\2\2\2\u016c\u016e\3\2\2\2\u016d\u016f\t\6\2\2\u016e\u016d\3\2"+
		"\2\2\u016e\u016f\3\2\2\2\u016f\u0170\3\2\2\2\u0170\u0172\7\61\2\2\u0171"+
		"\u0173\t\6\2\2\u0172\u0171\3\2\2\2\u0172\u0173\3\2\2\2\u0173\u0175\3\2"+
		"\2\2\u0174\u0176\t\5\2\2\u0175\u0174\3\2\2\2\u0176\u0177\3\2\2\2\u0177"+
		"\u0175\3\2\2\2\u0177\u0178\3\2\2\2\u0178\u0188\3\2\2\2\u0179\u017b\t\5"+
		"\2\2\u017a\u0179\3\2\2\2\u017b\u017c\3\2\2\2\u017c\u017a\3\2\2\2\u017c"+
		"\u017d\3\2\2\2\u017d\u0184\3\2\2\2\u017e\u0180\7\60\2\2\u017f\u0181\t"+
		"\5\2\2\u0180\u017f\3\2\2\2\u0181\u0182\3\2\2\2\u0182\u0180\3\2\2\2\u0182"+
		"\u0183\3\2\2\2\u0183\u0185\3\2\2\2\u0184\u017e\3\2\2\2\u0184\u0185\3\2"+
		"\2\2\u0185\u0186\3\2\2\2\u0186\u0188\7\'\2\2\u0187\u0169\3\2\2\2\u0187"+
		"\u017a\3\2\2\2\u0188N\3\2\2\2\u0189\u018a\7t\2\2\u018a\u018b\7g\2\2\u018b"+
		"\u018c\7o\2\2\u018c\u018d\7c\2\2\u018d\u018e\7k\2\2\u018e\u018f\7p\2\2"+
		"\u018f\u0190\7k\2\2\u0190\u0191\7p\2\2\u0191\u0192\7i\2\2\u0192P\3\2\2"+
		"\2\u0193\u0194\7m\2\2\u0194\u0195\7g\2\2\u0195\u0196\7r\2\2\u0196\u0197"+
		"\7v\2\2\u0197R\3\2\2\2\u0198\u0199\7d\2\2\u0199\u019a\7c\2\2\u019a\u019b"+
		"\7n\2\2\u019b\u019c\7c\2\2\u019c\u019d\7p\2\2\u019d\u019e\7e\2\2\u019e"+
		"\u019f\7g\2\2\u019fT\3\2\2\2\u01a0\u01a2\t\5\2\2\u01a1\u01a0\3\2\2\2\u01a2"+
		"\u01a3\3\2\2\2\u01a3\u01a1\3\2\2\2\u01a3\u01a4\3\2\2\2\u01a4V\3\2\2\2"+
		"\u01a5\u01a6\7\'\2\2\u01a6X\3\2\2\2\u01a7\u01a9\7&\2\2\u01a8\u01aa\t\7"+
		"\2\2\u01a9\u01a8\3\2\2\2\u01aa\u01ab\3\2\2\2\u01ab\u01a9\3\2\2\2\u01ab"+
		"\u01ac\3\2\2\2\u01ac\u01b0\3\2\2\2\u01ad\u01af\t\b\2\2\u01ae\u01ad\3\2"+
		"\2\2\u01af\u01b2\3\2\2\2\u01b0\u01ae\3\2\2\2\u01b0\u01b1\3\2\2\2\u01b1"+
		"Z\3\2\2\2\u01b2\u01b0\3\2\2\2\u01b3\u01b5\7B\2\2\u01b4\u01b6\t\t\2\2\u01b5"+
		"\u01b4\3\2\2\2\u01b6\u01b7\3\2\2\2\u01b7\u01b5\3\2\2\2\u01b7\u01b8\3\2"+
		"\2\2\u01b8\u01bc\3\2\2\2\u01b9\u01bb\t\n\2\2\u01ba\u01b9\3\2\2\2\u01bb"+
		"\u01be\3\2\2\2\u01bc\u01ba\3\2\2\2\u01bc\u01bd\3\2\2\2\u01bd\\\3\2\2\2"+
		"\u01be\u01bc\3\2\2\2\u01bf\u01c1\t\13\2\2\u01c0\u01bf\3\2\2\2\u01c1\u01c2"+
		"\3\2\2\2\u01c2\u01c0\3\2\2\2\u01c2\u01c3\3\2\2\2\u01c3^\3\2\2\2\27\2\u009c"+
		"\u00a1\u00aa\u00ac\u00ba\u0163\u016b\u016e\u0172\u0177\u017c\u0182\u0184"+
		"\u0187\u01a3\u01ab\u01b0\u01b7\u01bc\u01c2\3\b\2\2";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}