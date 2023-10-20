import pytest
from ccnp import Quote
from pytdxattest.tdreport  import TdReport

class Testquote:
    def test_integral_quote(self):
        quote1 = Quote.get_quote()
        quote2 = Quote.get_quote()
        assert quote1.quote_type==quote2.quote_type
        assert quote1.quote!=quote2.quote
    def test_quote_rtmrs(self):
        quote = Quote.get_quote()
        report=TdReport.get_td_report()
        report_rtmrs=report.td_info.rtmr_0+report.td_info.rtmr_1+report.td_info.rtmr_2+report.td_info.rtmr_3
        assert len(quote._rtmrs) == len (report_rtmrs)
        for i in range(0,len(quote._rtmrs)):
            assert quote._rtmrs[i] == report_rtmrs[i]

